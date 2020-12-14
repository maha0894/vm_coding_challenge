package models

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
	"vm_coding_challenge/db"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func setupTestDB() (sqlmock.Sqlmock, error) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		return mock, err
	}
	db.DB, err = gorm.Open("mysql", mockdb)
	return mock, err
}

func TestGetStats(t *testing.T) {
	t.Run("gets stats", func(t *testing.T) {
		now := time.Now()
		mock, err := setupTestDB()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		expStats := []HourlyStats{{ID: 1, CustomerID: 1, Time: now, RequestCount: 1, InvalidCount: 0}}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hourly_stats` WHERE (customer_id = ?)")).
			WithArgs(1).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "customer_id", "time", "request_count", "invalid_count"}).
				AddRow(1, 1, now, 1, 0))

		stats, _ := GetStats("customer", time.Time{}, 1)
		if !reflect.DeepEqual(expStats, stats) {
			t.Errorf("got %q, want %q", stats, expStats)
		}

		date := now.Truncate(24 * time.Hour)
		nextDay := date.AddDate(0, 0, 1).UTC()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hourly_stats` WHERE (time > ? AND time < ?)")).
			WithArgs(date, nextDay).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "customer_id", "time", "request_count", "invalid_count"}).
				AddRow(1, 1, now, 1, 0))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT sum(request_count)+sum(invalid_count) as total_day_requests FROM `hourly_stats` WHERE (time > ? AND time < ?)")).
			WithArgs(date, nextDay).
			WillReturnRows(sqlmock.
				NewRows([]string{"total_day_requests"}).
				AddRow(1))

		expRes := struct {
			Statistics       []HourlyStats `json:"statistics"`
			TotalDayRequests int           `json:"total_day_requests"`
		}{expStats, 1}

		stats, _ = GetStats("day", date, 0)
		if !reflect.DeepEqual(expRes, stats) {
			t.Errorf("got %q, want %q", stats, expRes)
		}
	})
}

func TestCheckValidity(t *testing.T) {
	t.Run("checks validity", func(t *testing.T) {
		req := Request{CustomerID: 0, TagID: 0, UserID: "", RemoteIP: "", Timestamp: 0}
		userAgent := ""
		err := CheckValidity(req, userAgent)
		if err == errors.New("customer ID not specified") {
			t.Errorf("got %q, want %q", err, errors.New("customer ID not specified"))
		}
		req.CustomerID = 1
		err = CheckValidity(req, userAgent)
		if err == errors.New("tag ID not specified") {
			t.Errorf("got %q, want %q", err, errors.New("tag ID not specified"))
		}
		req.TagID = 1
		err = CheckValidity(req, userAgent)
		if err == errors.New("user ID not specified") {
			t.Errorf("got %q, want %q", err, errors.New("user ID not specified"))
		}

		req.UserID = "user_id"
		err = CheckValidity(req, userAgent)
		if err == errors.New("remote IP not specified") {
			t.Errorf("got %q, want %q", err, errors.New("remote IP not specified"))
		}
		req.RemoteIP = "123.123.13.12"
		err = CheckValidity(req, userAgent)
		if err == errors.New("timestamp not specified") {
			t.Errorf("got %q, want %q", err, errors.New("timestamp not specified"))
		}
		req.Timestamp = 1500110000

		mock, err := setupTestDB()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `ua_blacklist` WHERE (ua = ?)")).
			WithArgs(userAgent).
			WillReturnError(gorm.ErrRecordNotFound)

		remIP, _ := strconv.Atoi(strings.ReplaceAll(req.RemoteIP, ".", ""))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `ip_blacklist` WHERE (ip = ?)")).
			WithArgs(remIP).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE `customer`.`id` = ? AND ((id = ?))")).
			WithArgs(req.CustomerID, req.CustomerID).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "name", "active"}).
				AddRow(req.CustomerID, "name", true))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `hourly_stats` WHERE (customer_id = ? AND time = ?)")).
			WithArgs(req.CustomerID, time.Unix(int64(req.Timestamp), 0).Truncate(time.Hour)).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `hourly_stats`").
			WithArgs(req.CustomerID, time.Unix(int64(req.Timestamp), 0).Truncate(time.Hour), 1, 0).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = CheckValidity(req, userAgent)
		if err != nil {
			t.Errorf("got %v, want %v", err, nil)
		}
	})
}
