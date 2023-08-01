package users

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/users-service/internal/domain"
	"github.com/olad5/productive-pulse/users-service/internal/infra"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo       infra.UserRepository
	configurations *config.Configurations
}

var (
	ErrUserAlreadyExists = "email already exist"
	ErrPasswordIncorrect = "incorrect credentials"
	ErrUserNotFound      = "user not found"
	ErrInvalidToken      = errors.New("invalid token")
)

func NewUserService(userRepo infra.UserRepository, configurations *config.Configurations) (*UserService, error) {
	if userRepo == nil {
		return &UserService{}, errors.New("UserService failed to initialize")
	}
	return &UserService{userRepo, configurations}, nil
}

func (u *UserService) CreateUser(ctx context.Context, tracer trace.Tracer, firstName, lastName, email, password string) (domain.User, error) {
	ctx, span := tracer.Start(ctx, "CreateUser")
	defer span.End()
	existingUser, err := u.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser.Email == email {
		return domain.User{}, errors.New(ErrUserAlreadyExists)
	}
	if err != nil {
		return domain.User{}, err
	}

	hashedPassword, err := hashAndSalt([]byte(password))
	newUser := domain.User{
		ID:        uuid.New(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  hashedPassword,
	}
	if err != nil {
		return domain.User{}, err
	}

	err = u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, err
	}
	return newUser, nil
}

func (u *UserService) LogUserIn(ctx context.Context, tracer trace.Tracer, email, password string) (string, error) {
	ctx, span := tracer.Start(ctx, "LogUserIn")
	defer span.End()
	existingUser, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil && err.Error() == infra.ErrRecordNotFound {
		return "", errors.New(ErrUserNotFound)
	}
	if isPasswordCorrect := comparePasswords(existingUser.Password, []byte(password)); isPasswordCorrect == false {
		return "", errors.New(ErrPasswordIncorrect)
	}

	accessToken, err := generateJWT(existingUser, u.configurations.UserServiceSecretKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (u *UserService) VerifyUser(ctx context.Context, tracer trace.Tracer, authHeader string) (string, error) {
	ctx, span := tracer.Start(ctx, "VerifyUser")
	defer span.End()
	const Bearer = "Bearer "
	if authHeader != "" && strings.HasPrefix(authHeader, Bearer) {
		token := strings.TrimPrefix(authHeader, Bearer)
		return verifyJWT(token, u.configurations.UserServiceSecretKey)
	}

	return "", ErrInvalidToken
}

func generateJWT(user domain.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("Error generating JWT token")
	}
	return tokenString, nil
}

func verifyJWT(tokenString, secret string) (string, error) {
	if tokenString == "" {
		return "", ErrInvalidToken
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", errors.New("error decoding jwt")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return "", errors.New("expired token")
		}
		userId, ok := claims["sub"]
		if ok == true && userId != nil {
			return userId.(string), nil
		}
		return "", ErrInvalidToken
	} else {
		return "", ErrInvalidToken
	}
}

func hashAndSalt(plainPassword []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(plainPassword, bcrypt.MinCost)
	if err != nil {
		return "", errors.New("error hashing password")
	}
	return string(hash), nil
}

func comparePasswords(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		return false
	}

	return true
}
