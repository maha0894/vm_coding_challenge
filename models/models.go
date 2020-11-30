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

// Validate is the function which checks Request fields
func (r *Request) Validate() (err error) {
	switch {
	case r.CustomerID == 0:
		return errors.New("customer ID not specified")
	case r.TagID == 0:
		return errors.New("tag ID not specified")
	case r.UserID == "":
		return errors.New("user ID not specified")
	case r.RemoteIP == "":
		return errors.New("remote IP not specified")
	case r.Timestamp == 0:
		return errors.New("timestamp not specified")
	}
	return err
}

// CheckValidity is the function which checks validity of Request
func CheckValidity(req Request) (err error) {
	err = req.Validate()
	if err == nil {
		// TODO check remote IP address which is in the blacklist
		// TODO check user agent
		// TODO check customer ID not found in the database or disabled
	}
	if req.CustomerID != 0 {
		// count stats
	}
	return
}

// GetStats is the function which returns stats based on params
func GetStats(by string, timestamp, customerID int) (stats Stats, err error) {

	return
}
