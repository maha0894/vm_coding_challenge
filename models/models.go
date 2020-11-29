package models

import "errors"

type Request struct {
	CustomerID int    `json:"customerID"`
	TagID      int    `json:"tagID"`
	UserID     string `json:"userID"`
	RemoteIP   string `json:"remoteIP"`
	Timestamp  int    `json:"timestamp"`
}

type Stats struct {
	ID           int `json:"id"`
	CustomerID   int `json:"customer_id"`
	Time         int `json:"time"`
	RequestCount int `json:"request_count"`
	InvalidCount int `json:"invalid_count"`
}

// Validate is the function which check incoming request
func (r *Request) Validate() (err error) {
	// TODO check remote IP address which is in the blacklist
	// TODO check user agent
	switch {
	case r.CustomerID == 0:
		return errors.New("customer ID not specified")
	}
	// TODO check customer ID not found in the database or disabled
	return err
}

// GetStats is the function which returns stats based on params
func GetStats(by string, timestamp, customerID int) (stats Stats, err error) {

	return
}
