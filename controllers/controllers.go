package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"vm_coding_challenge/models"

	"github.com/gorilla/mux"

	useragent "github.com/mileusna/useragent"
)

// Request handles /request route
func Request(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		JSONResponse(w, "Invalid JSON structure", http.StatusBadRequest)
		return
	}
	ua := useragent.Parse(r.Header.Get("User-Agent"))
	err = models.CheckValidity(req, ua.Name)
	if err != nil {
		JSONResponse(w, "Request rejected "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, "Request accepted", http.StatusOK)
}

// Statistics handles /stats route
func Statistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date, _ := time.Parse("2006-01-02", vars["day"])
	customerID, _ := strconv.Atoi(vars["id"])
	stats, err := models.GetStats(vars["by"], date, customerID)
	if err != nil {
		fmt.Println("Error getting statistics: ", err)
		JSONResponse(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	JSONResponse(w, stats, http.StatusOK)
}

// JSONResponse attempts to set the status code, c, and marshal the given interface, d, into a response that
// is written to the given ResponseWriter.
func JSONResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		fmt.Println("Error creating JSON response: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}
