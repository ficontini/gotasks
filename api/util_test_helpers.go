package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func makeRequest(method, path, token string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", token))
	return req
}
func makeUnauthenticatedRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	req.Header.Add("Content-Type", "application/json")
	return req
}
func testRequest(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
func marshallParamsToJSON(t *testing.T, params interface{}) []byte {
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func checkStatusCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("expected %d status code, but got %d", expected, actual)
	}
}
