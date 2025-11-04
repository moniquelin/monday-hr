package main

import (
	"encoding/json" // New import
	"net/http"
)

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which holds the information that we want to send in the response
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	// Pass the map to the json.Marshal() function
	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
