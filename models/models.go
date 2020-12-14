package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"vm_coding_challenge/db"

	"github.com/jinzhu/gorm"
)

type Request struct {
	CustomerID int    `json:"customerID"`
	TagID      int    `json:"tagID"`
	UserID     string `json:"userID"`
	RemoteIP   string `json:"remoteIP"`
	Timestamp  int    `json:"timestamp"`
}

type HourlyStats struct {
	ID           int       `json:"id"`
	CustomerID   int       `json:"customer_id"`
	Time         time.Time `json:"time"`
	RequestCount int       `json:"request_count"`
	InvalidCount int       `json:"invalid_count"`
}

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type IPBlacklist struct {
	IP int `json:"ip"`
}

type UABlacklist struct {
	UA string `json:"ua"`
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
func CheckValidity(req Request, userAgent string) (err error) {
	var valid bool
	err = req.Validate()
	if err == nil {
		// check user agent
		ua := UABlacklist{UA: userAgent}
		if ua.valid() {
			// check remote IP address which is in the blacklist
			remIP, _ := strconv.Atoi(strings.ReplaceAll(req.RemoteIP, ".", ""))
			ip := IPBlacklist{IP: remIP}
			if ip.valid() {
				// check customer ID not found in the database or disabled
				c := Customer{ID: req.CustomerID}
				if c.valid() {
					valid = true
				} else {
					// customer not found
					return errors.New("customer not found")
				}
			}
		}
	} else {
		// missing field(s)
		return err
	}
	if req.CustomerID != 0 {
		// count stats
		stats := HourlyStats{CustomerID: req.CustomerID, Time: time.Unix(int64(req.Timestamp), 0).Truncate(time.Hour)}
		updateStats(stats, valid)
		if valid {
			req.handleValidRequest()
		}
	}
	return
}

func (ua *UABlacklist) valid() bool {
	err := db.DB.Table("ua_blacklist").Where("ua = ?", ua.UA).Find(&ua).Error
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func (ip *IPBlacklist) valid() bool {
	err := db.DB.Table("ip_blacklist").Where("ip = ?", ip.IP).Find(&ip).Error
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func (c *Customer) valid() bool {
	err := db.DB.Table("customer").Where("id = ?", c.ID).Find(&c).Error
	if err != nil || !c.Active {
		return false
	}
	return true
}

func updateStats(stats HourlyStats, valid bool) (err error) {
	err = db.DB.Table("hourly_stats").
		Where("customer_id = ? AND time = ?", stats.CustomerID, stats.Time).Find(&stats).Error
	if err == gorm.ErrRecordNotFound {
		// create
		if valid {
			stats.RequestCount = 1
		} else {
			stats.InvalidCount = 1
		}
		err = db.DB.Table("hourly_stats").Create(&stats).Error
		return err
	} else if err != nil {
		return err
	}
	// update
	if valid {
		err = db.DB.Table("hourly_stats").
			Where("id = ?", stats.ID).
			Update("request_count", stats.RequestCount+1).Error
	} else {
		err = db.DB.Table("hourly_stats").
			Where("id = ?", stats.ID).
			Update("invalid_count", stats.InvalidCount+1).Error
	}
	return err
}

func (r *Request) handleValidRequest() {

}

// GetStats is the function which returns stats based on params
func GetStats(by string, date time.Time, customerID int) (res interface{}, err error) {
	switch by {
	case "customer":
		var stats []HourlyStats
		err = db.DB.Table("hourly_stats").
			Where("customer_id = ?", customerID).Find(&stats).Error
		return stats, err
	case "day":
		stats := struct {
			Statistics       []HourlyStats `json:"statistics"`
			TotalDayRequests int           `json:"total_day_requests"`
		}{}
		nextDay := date.AddDate(0, 0, 1).UTC()
		err = db.DB.Table("hourly_stats").
			Where("time > ? AND time < ?", date, nextDay).Find(&stats.Statistics).Error
		if err != nil {
			return
		}
		// count total
		err = db.DB.Table("hourly_stats").
			Where("time > ? AND time < ?", date, nextDay).
			Select("sum(request_count)+sum(invalid_count) as total_day_requests").
			Find(&stats).Error
		return stats, err
	}
	return
}
