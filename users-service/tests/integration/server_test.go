//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/users-service/internal/app/router"
	"github.com/olad5/productive-pulse/users-service/internal/app/server"
	"github.com/olad5/productive-pulse/users-service/internal/handlers"
	"github.com/olad5/productive-pulse/users-service/internal/infra/postgres"
	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
)

var (
	appRouter http.Handler
	svr       *server.Server
)

func TestMain(m *testing.M) {
	configurations := config.GetConfig("../config/.test.env")
	ctx := context.Background()

	userRepo, err := postgres.NewPostgresRepo(ctx, configurations.UserServiceDBUrl)
	if err != nil {
		log.Fatal("Error Initializing User Repo")
	}

	err = userRepo.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	userService, err := users.NewUserService(userRepo, configurations)
	if err != nil {
		log.Fatal("Error Initializing UserService")
	}

	userHandler, err := handlers.NewHandler(*userService)
	if err != nil {
		log.Fatal("failed to create the User handler: ", err)
	}
	appRouter = router.NewHttpRouter(*userHandler)
	svr = server.CreateNewServer(appRouter)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func executeRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func TestRegister(t *testing.T) {
	route := "/users"
	t.Run("test for invalid json request body",
		func(t *testing.T) {
			req, _ := http.NewRequest("POST", route, nil)
			response := executeRequest(req, svr)
			assertStatusCode(t, http.StatusBadRequest, response.Code)
		},
	)
	t.Run(`Given a valid user registration request, when the user submits the request, 
    then the server should respond with a success status code, and the user's account 
    should be created in the database.`,
		func(t *testing.T) {
			email := "will" + fmt.Sprint(generateUniqueId()) + "@gmail.com"
			requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "first_name": "will",
      "last_name": "hansen",
      "password": "some-random-password"
      }`, email))
			req, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))
			response := executeRequest(req, svr)
			assertStatusCode(t, http.StatusOK, response.Code)
			data := parseResponse(response)["data"].(map[string]interface{})
			assertResponseMessage(t, data["email"].(string), email)
		},
	)

	t.Run(`Given a user registration request with an email address that already exists,
    when the user submits the request, then the server should respond with an error
    status code, and the response should indicate that the email address is already
    taken. `,
		func(t *testing.T) {
			email := "will@gmail.com"
			requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "first_name": "will",
      "last_name": "hansen",
      "password": "some-random-password"
      }`, email))
			req, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))
			response := executeRequest(req, svr)

			secondRequestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "first_name": "will",
      "last_name": "hansen",
      "password": "passcode"
      }`, email))
			req, _ = http.NewRequest("POST", route, bytes.NewBuffer(secondRequestBody))

			assertStatusCode(t, http.StatusBadRequest, response.Code)
			message := parseResponse(response)["message"].(string)
			assertResponseMessage(t, message, "email already exist")
		},
	)
}

func TestLogin(t *testing.T) {
	route := "/users/login"
	t.Run("test for invalid json request body",
		func(t *testing.T) {
			req, _ := http.NewRequest("POST", route, nil)
			response := executeRequest(req, svr)
			assertStatusCode(t, http.StatusBadRequest, response.Code)
		},
	)
	t.Run(`Given a user attempts to log in with valid credentials,
    when they make a POST request to the login endpoint with their username and password,
    then they should receive a 200 OK response,
    and the response should contain a JSON web token (JWT) in the 'token' field,
    and the token should be valid and properly signed.`,
		func(t *testing.T) {
			email := "will@gmail.com"
			requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "password": "some-random-password"
      }`, email))
			req, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))
			response := executeRequest(req, svr)

			assertStatusCode(t, http.StatusOK, response.Code)
			responseBody := parseResponse(response)
			message := responseBody["message"].(string)
			assertResponseMessage(t, message, "user logged in successfully")

			data := responseBody["data"].(map[string]interface{})

			accessToken, exists := data["access_token"]
			if !exists {
				t.Error("Missing 'accesstoken' key in the JSON response")
			}

			_, isString := accessToken.(string)
			if !isString {
				t.Error("'accesstoken' value is not a string")
			}
		},
	)

	t.Run(`Given a user tries to log in with an account that does not exist,
    when they make a POST request to the login endpoint with a non-existent email,
    then they should receive a 404 Not Found response,
    and the response should contain an error message indicating that the account 
    does not exist.`,
		func(t *testing.T) {
			email := "emailnoexist@gmail.com"
			requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "first_name": "mike",
      "last_name": "wilson",
      "password": "some-random-password"
      }`, email))
			req, _ := http.NewRequest("POST", route, bytes.NewBuffer(requestBody))
			response := executeRequest(req, svr)

			assertStatusCode(t, http.StatusBadRequest, response.Code)
			message := parseResponse(response)["message"].(string)
			assertResponseMessage(t, message, "incorrect credentials")
		},
	)
}

func TestAuth(t *testing.T) {
	route := "/users/auth"
	t.Run(` Given a user trying to access a protected endpoint without providing a JWT token
      When they make a GET request to the protected endpoint without an authorization header
      Then they should receive a 401 Unauthorized response
      And the response should contain an error message indicating that a token is required for access `,
		func(t *testing.T) {
			req, _ := http.NewRequest("GET", route, nil)

			response := executeRequest(req, svr)
			assertStatusCode(t, http.StatusUnauthorized, response.Code)
		},
	)
	t.Run(`Given a user with valid Bearer Token
      When they make a GET request to the auth endpoint 
      Then they should receive a 200 OK response
      And the response should contain a JSON with the 'user_id' field `,
		func(t *testing.T) {
			email := "will@gmail.com"
			requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "password": "some-random-password"
      }`, email))
			loginReq, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(requestBody))
			loginResponse := executeRequest(loginReq, svr)
			loginResponseBody := parseResponse(loginResponse)
			loginData := loginResponseBody["data"].(map[string]interface{})
			accessToken := loginData["access_token"]
			token := accessToken.(string)

			req, _ := http.NewRequest("GET", route, nil)
			req.Header.Set("Authorization", "Bearer "+token)

			response := executeRequest(req, svr)
			assertStatusCode(t, http.StatusOK, response.Code)
		},
	)
	t.Run(`Given a user with an invalid JWT token signature
    When they make a GET request to a protected endpoint with the invalid token in the authorization header
    Then they should receive a 401 Unauthorized response
    And the response should contain an error message indicating that the token is not valid `,
		func(t *testing.T) {
			req, _ := http.NewRequest("GET", route, nil)
			token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOGFhNzQwMzEtMzkzZC00MDFlLWE1ZDMtYmE3MjA4OWFiZTQwIn0.fbQDiaMeWYt_nAWsIbxqmwSgQKOYYz1OtlJikAHMqVE"
			req.Header.Set("Authorization", "Bearer "+token)

			response := executeRequest(req, svr)
			assertStatusCode(t, http.StatusUnauthorized, response.Code)
		},
	)
}

func assertStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func assertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("got message: %q expected: %q", got, expected)
	}
}

func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func generateUniqueId() int {
	MAX_INT := 7935425686241
	b := new(big.Int).SetInt64(int64(MAX_INT))
	randomBigInt, _ := rand.Int(rand.Reader, b)
	randomeNewInt := int(randomBigInt.Int64())
	return randomeNewInt
}
