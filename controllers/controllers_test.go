package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
	"time"
	"vm_coding_challenge/db"
	"vm_coding_challenge/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func setupTestDB() (sqlmock.Sqlmock, error) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		return mock, err
	}
	db.DB, err = gorm.Open("mysql", mockdb)
	return mock, err
}

func TestStatistics(t *testing.T) {
	t.Run("returns stats", func(t *testing.T) {
		now := time.Now().UTC()

		mock, err := setupTestDB()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hourly_stats` WHERE (customer_id = ?)")).
			WithArgs(1).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "customer_id", "time", "request_count", "invalid_count"}).
				AddRow(1, 1, now, 1, 0))

		request, _ := http.NewRequest(http.MethodGet, "/stats/customer/1", nil)
		request = mux.SetURLVars(request, map[string]string{
			"by": "customer",
			"id": "1",
		})
		response := httptest.NewRecorder()

		Statistics(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		var got []models.HourlyStats
		err = json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q, '%v'", response.Body, err)
		}

		want := []models.HourlyStats{{ID: 1, CustomerID: 1, Time: now, RequestCount: 1, InvalidCount: 0}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

	})
}

func TestRequest(t *testing.T) {
	t.Run("expects request", func(t *testing.T) {
		body := []byte(`{"customerID":0,"tagID":2,"userID":"aaaaaaaa-bbbb-cccc-1111-222222222222","remoteIP":"123.234.56.78","timestamp":1500000000}`)

		request, _ := http.NewRequest(http.MethodPut, "/request", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		Request(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)

		resBody, _ := ioutil.ReadAll(response.Body)
		var got string
		json.Unmarshal(resBody, &got)

		want := "Request rejected customer ID not specified"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}

	})
}
