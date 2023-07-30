package tests

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olad5/productive-pulse/pkg/app/server"
)

func ExecuteRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func AssertStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("got message: %q expected: %q", got, expected)
	}
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func GenerateUniqueId() int {
	MAX_INT := 7935425686241
	b := new(big.Int).SetInt64(int64(MAX_INT))
	randomBigInt, _ := rand.Int(rand.Reader, b)
	randomeNewInt := int(randomBigInt.Int64())
	return randomeNewInt
}
