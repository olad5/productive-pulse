//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/pkg/app/server"
	tests "github.com/olad5/productive-pulse/pkg/tests"
	"github.com/olad5/productive-pulse/todo-service/internal/app/router"
	"github.com/olad5/productive-pulse/todo-service/internal/handlers"
	"github.com/olad5/productive-pulse/todo-service/internal/infra/mongo"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
)

var svr *server.Server

func TestMain(m *testing.M) {
	configurations := config.GetConfig("../config/.test.env")
	ctx := context.Background()

	todoRepo, err := mongo.NewMongoRepo(ctx, configurations.TodoServiceDBConnectionString)
	if err != nil {
		log.Fatal("Error Initializing Todo Repo")
	}

	todoService, err := todos.NewTodoService(todoRepo, configurations)
	if err != nil {
		log.Fatal("Error Initializing TodoService")
	}
	userService := &StubUserService{
		client: &http.Client{},
		url:    "",
	}

	todoHandler, err := handlers.NewTodoHandler(*todoService, userService)
	if err != nil {
		log.Fatal("failed to create the Todo handler: ", err)
	}
	appRouter := router.NewHttpRouter(*todoHandler)
	svr = server.CreateNewServer(appRouter)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCreateTodo(t *testing.T) {
	route := "/todos"
	t.Run("test for invalid json request body",
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, route, nil)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusBadRequest, response.Code)
		},
	)
	t.Run(`Given an authenticated user with valid credentials and a valid JWT token
    When they make a POST request to the create todo endpoint with a valid todo object
    Then they should receive a 200 Created response
    And the response should contain the created todo object with a unique ID
    And the todo object should have the correct text
    `,
		func(t *testing.T) {
			text := "some random text with an id :" + fmt.Sprint(tests.GenerateUniqueId())
			requestBody := []byte(fmt.Sprintf(`{
      "text": "%s"
      }`, text))
			req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			data := tests.ParseResponse(response)["data"].(map[string]interface{})
			tests.AssertResponseMessage(t, data["text"].(string), text)
		},
	)
	t.Run(`Given an unauthenticated user or a user with an invalid JWT token
        When they make a POST request to the create todo endpoint
        Then they should receive a 401 Unauthorized response
        And the response should contain an error message indicating that 
        authentication is required
    `,
		func(t *testing.T) {
			text := "some random text with an id :" + fmt.Sprint(tests.GenerateUniqueId())
			requestBody := []byte(fmt.Sprintf(`{
      "text": "%s"
      }`, text))
			req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
			req.Header.Set("Authorization", "Bearer "+inValidToken)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusUnauthorized, response.Code)
		},
	)
}

func TestGetTodo(t *testing.T) {
	route := "/todos"
	t.Run(` Given an authenticated user with valid credentials and a valid JWT token
      And there exists a todo with a specific ID in the database
      When they make a GET request to the get todo endpoint with the valid todo ID
      Then they should receive a 200 OK response
      And the response should contain the requested todo object with the correct ID,
      and text
    `,
		func(t *testing.T) {
			text := "some random text with an id :" + fmt.Sprint(tests.GenerateUniqueId())
			createRequestBody := []byte(fmt.Sprintf(`{
			"text": "%s"
			}`, text))
			createReq, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(createRequestBody))
			createReq.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			createResponse := tests.ExecuteRequest(createReq, svr)
			createData := tests.ParseResponse(createResponse)["data"].(map[string]interface{})
			id := createData["id"].(string)

			req, _ := http.NewRequest(http.MethodGet, route+"/"+id, nil)
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			data := tests.ParseResponse(response)["data"].(map[string]interface{})
			tests.AssertResponseMessage(t, data["id"].(string), id)
		},
	)

	t.Run(` Given an authenticated user with valid credentials and a valid JWT token
      And there is no todo with the specified ID in the database
      When they make a GET request to the get todo endpoint with an invalid or non-existent todo ID
      Then they should receive a 404 Not Found response
      And the response should contain an error message indicating that the requested todo was not found
    `,
		func(t *testing.T) {
			id := "caf3d5c8-0db7-4b27-a02e-e5a664684568"
			req, _ := http.NewRequest(http.MethodGet, route+"/"+id, nil)
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusNotFound, response.Code)
		},
	)
	t.Run(`Given an authenticated user with valid credentials and a valid JWT token
    When they make a GET request to the get todo endpoint with an invalid todo ID format
    Then they should receive a 400 Bad Request response
    And the response should contain an error message indicating that the todo ID is invalid
    `,
		func(t *testing.T) {
			id := "caf3d5c8-0db7-4b27-a02e-e5a6646868"
			req, _ := http.NewRequest(http.MethodGet, route+"/"+id, nil)
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusBadRequest, response.Code)
		},
	)

	t.Run(` Given an authenticated user with valid credentials and a valid JWT token
      And there exists a todo with a specific ID in the database
      When they make a GET request to the get todo endpoint with the valid todo ID
      And the user did not create the todo.
      Then they should receive a 401 Unauthorized response And the response should
      contain an error message indicating that authentication is required
    `,
		func(t *testing.T) {
			text := "some random text with an id :" + fmt.Sprint(tests.GenerateUniqueId())
			createRequestBody := []byte(fmt.Sprintf(`{
			"text": "%s"
			}`, text))
			createReq, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(createRequestBody))
			createReq.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			createResponse := tests.ExecuteRequest(createReq, svr)
			createData := tests.ParseResponse(createResponse)["data"].(map[string]interface{})
			id := createData["id"].(string)
			req, _ := http.NewRequest(http.MethodGet, route+"/"+id, nil)
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser2)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusUnauthorized, response.Code)
		},
	)
}

func TestGetTodos(t *testing.T) {
	route := "/todos"
	t.Run(` Given an authenticated user with valid credentials and a valid JWT token
      And there exists a todo with a specific ID in the database
      When they make a GET request to the get todo endpoint with the valid todo ID
      Then they should receive a 200 OK response
      And the response should contain the requested todo object with the correct ID,
      and text
    `,
		func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, route, nil)
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			data := tests.ParseResponse(response)["data"].([]interface{})
			if len(data) > 1 != true {
				t.Errorf("user should have more than 1 todo")
			}
		},
	)
}

func TestUpdateTodo(t *testing.T) {
	route := "/todos"
	t.Run(` Given an authenticated user with valid credentials and a valid JWT token
      And there exists a todo with a specific ID in the database
      When they make a GET request to the get todo endpoint with the valid todo ID
      Then they should receive a 200 OK response
      And the response should contain the requested todo object with the correct ID,
      and text
    `,
		func(t *testing.T) {
			text := "some random text with an id :" + fmt.Sprint(tests.GenerateUniqueId())
			createRequestBody := []byte(fmt.Sprintf(`{
			"text": "%s"
			}`, text))
			createReq, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(createRequestBody))
			createReq.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			createResponse := tests.ExecuteRequest(createReq, svr)
			createData := tests.ParseResponse(createResponse)["data"].(map[string]interface{})
			id := createData["id"].(string)

			updatedText := "updated text with an id: " + fmt.Sprint(tests.GenerateUniqueId())
			updatedRequestBody := []byte(fmt.Sprintf(`{
			"text": "%s"
			}`, updatedText))

			req, _ := http.NewRequest(http.MethodPatch, route+"/"+id, bytes.NewBuffer(updatedRequestBody))
			req.Header.Set("Authorization", "Bearer "+ValidTokenForUser1)
			response := tests.ExecuteRequest(req, svr)
			tests.AssertStatusCode(t, http.StatusOK, response.Code)
			data := tests.ParseResponse(response)["data"].(map[string]interface{})
			tests.AssertResponseMessage(t, data["text"].(string), updatedText)
		},
	)
}
