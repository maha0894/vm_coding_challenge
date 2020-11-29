package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func TestStatistics(t *testing.T) {
	t.Run("returns stats", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/stats", nil)
		response := httptest.NewRecorder()

		Statistics(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		// TODO check body

	})
}

func TestRequest(t *testing.T) {
	t.Run("expects request", func(t *testing.T) {
		body := []byte(`{"customerID":1,"tagID":2,"userID":"aaaaaaaa-bbbb-cccc-1111-222222222222","remoteIP":"123.234.56.78","timestamp":1500000000}`)

		request, _ := http.NewRequest(http.MethodPut, "/request", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		Request(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		// TODO check body

	})
}
