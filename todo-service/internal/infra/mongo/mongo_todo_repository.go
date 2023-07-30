package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/olad5/productive-pulse/todo-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	todos *mongo.Collection
}

var contextTimeoutDuration = 5 * time.Second

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = byte(0x04)

	mongoRegistry = bson.NewRegistryBuilder().
			RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue)).
			RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue)).
			Build()
)

func NewMongoRepo(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString).SetRegistry(mongoRegistry))
	if err != nil {
		return nil, fmt.Errorf("failed to create a mongo client: %w", err)
	}

	todoCollection := client.Database("todo-service").Collection("todos")

	return &MongoRepository{
		todos: todoCollection,
	}, nil
}

func (m *MongoRepository) CreateTodo(ctx context.Context, todo domain.Todo) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoTodo := toMongoTodo(todo)
	_, err := m.todos.InsertOne(ctx, mongoTodo)
	if err != nil {
		return fmt.Errorf("failed to persist todo: %w", err)
	}
	return nil
}

func (m *MongoRepository) GetTodo(ctx context.Context, userId, todoId uuid.UUID) (domain.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	todo := mongoTodo{}
	err := m.todos.FindOne(ctx, bson.M{"_id": todoId}).Decode(&todo)
	if err != nil {
		return domain.Todo{}, errors.New("record not found")
	}
	domainTodo := toTodo(todo)
	if domainTodo.UserId != userId {
		return domain.Todo{}, errors.New("current user is not owner of this todo")
	}

	return domainTodo, nil
}

func (m *MongoRepository) GetTodos(ctx context.Context, userId uuid.UUID) ([]domain.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	filter := bson.M{"user_id": userId}

	cursor, err := m.todos.Find(ctx, filter)
	if err != nil {
		return []domain.Todo{}, errors.New("errors getting todos")
	}
	defer cursor.Close(ctx)
	var mongoTodos []mongoTodo
	if err = cursor.All(ctx, &mongoTodos); err != nil {
		return []domain.Todo{}, errors.New("errors getting todos")
	}
	var domainTodos []domain.Todo
	for _, mongoTodo := range mongoTodos {
		domainTodos = append(domainTodos, toTodo(mongoTodo))
	}

	return domainTodos, nil
}

func (m *MongoRepository) UpdateTodo(ctx context.Context, todo domain.Todo) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutDuration)
	defer cancel()

	mongoTodo := toMongoTodo(todo)
	filter := bson.M{"_id": todo.ID}
	updatedDoc := bson.M{
		"$set": mongoTodo,
	}
	_, err := m.todos.UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		return fmt.Errorf("failed to persist todo: %w", err)
	}
	return nil
}

func (m *MongoRepository) Ping(ctx context.Context) error {
	if _, err := m.todos.EstimatedDocumentCount(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}
	return nil
}

type mongoTodo struct {
	ID        uuid.UUID `bson:"_id"`
	UserId    uuid.UUID `bson:"user_id"`
	Text      string    `bson:"text"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func toMongoTodo(todo domain.Todo) mongoTodo {
	return mongoTodo{
		ID:        todo.ID,
		UserId:    todo.UserId,
		Text:      todo.Text,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func toTodo(m mongoTodo) domain.Todo {
	return domain.Todo{
		ID:        m.ID,
		UserId:    m.UserId,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// NOTE: how to store uuid in mongodb
// https://stackoverflow.com/questions/64723089/how-to-store-a-uuid-in-mongodb-with-golang
// https://gist.github.com/SupaHam/3afe982dc75039356723600ccc91ff77
func uuidEncodeValue(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}
