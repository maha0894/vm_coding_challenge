package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vm_coding_challenge/models"

	"github.com/gorilla/mux"
)

// Request handles /request route
func Request(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		JSONResponse(w, "Invalid JSON structure", http.StatusBadRequest)
		return
	}
	err = models.CheckValidity(req)
	if err != nil {
		JSONResponse(w, "Request rejected "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, "Request accepted", http.StatusOK)
}

// Statistics handles /stats route
func Statistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	timestamp, _ := strconv.Atoi(vars["day"])
	customerID, _ := strconv.Atoi(vars["id"])
	stats, err := models.GetStats(vars["by"], timestamp, customerID)
	if err != nil {
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
