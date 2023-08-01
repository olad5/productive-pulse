package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jaswdr/faker"
)

var proxyBaseURL = "http://localhost"

var (
	accessToken    string
	existingTodoId string
)

func main() {
	newUserEmail, newUserPassword := createUser()
	accessToken = logInUser(newUserEmail, newUserPassword)

	// success request
	existingTodoId = createTodo(accessToken, faker.New().Lorem().Text(10))
	time.Sleep(500 * time.Millisecond)

	// failing request
	createTodo("", faker.New().Lorem().Text(10))
	time.Sleep(500 * time.Millisecond)

	// success request
	getTodo(accessToken, existingTodoId)
	time.Sleep(500 * time.Millisecond)

	// failing request
	getTodo(accessToken, "1eec0938-50af-4faa-b7-fb409954f4f6")
	time.Sleep(500 * time.Millisecond)
}

func createUser() (string, string) {
	route := proxyBaseURL + "/users"
	email := faker.New().Internet().Email()
	first_name := faker.New().Person().FirstName()
	last_name := faker.New().Person().LastName()
	password := "some-random-password"
	requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "first_name": "%s",
      "last_name": "%s",
      "password": "%s"
      }`, email, first_name, last_name, password))
	req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
	res := newMakeRequest(req)
	defer res.Body.Close()
	return email, password
}

func logInUser(email, password string) string {
	route := proxyBaseURL + "/users/login"
	requestBody := []byte(fmt.Sprintf(`{
      "email": "%s",
      "password": "%s"
      }`, email, password))
	req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
	res := newMakeRequest(req)
	defer res.Body.Close()

	responseBody := parseResponse(res)
	bs, _ := json.Marshal(responseBody)
	fmt.Println(string(bs))
	if res.StatusCode == http.StatusOK {
		data := responseBody["data"].(map[string]interface{})
		accessToken := data["access_token"].(string)
		return accessToken
	}
	return ""
}

func createTodo(accessToken, text string) string {
	route := proxyBaseURL + "/todos"
	requestBody := []byte(fmt.Sprintf(`{
      "text": "%s"
      }`, text))
	req, _ := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res := newMakeRequest(req)
	defer res.Body.Close()
	responseBody := parseResponse(res)
	bs, _ := json.Marshal(responseBody)
	fmt.Println(string(bs))
	if res.StatusCode == http.StatusOK {
		data := responseBody["data"].(map[string]interface{})
		todoId := data["id"].(string)
		return todoId
	}
	return ""
}

func getTodo(accessToken, todoId string) {
	route := proxyBaseURL + "/todos" + "/" + todoId
	req, _ := http.NewRequest(http.MethodGet, route, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res := newMakeRequest(req)
	defer res.Body.Close()
}

func newMakeRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	url := req.URL.Path
	if res.StatusCode == http.StatusOK {
		fmt.Printf("Request to %s succeeded with status code %d\n", url, res.StatusCode)
	} else {
		fmt.Printf("Request to %s failed with status code %d\n", url, res.StatusCode)
	}
	return res
}

func parseResponse(w *http.Response) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}
