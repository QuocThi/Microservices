package main

import (
	"net/http"

	"log-service/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.HandleError(w, "marshal request failed", err)
		return
	}

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.DB.Insert(event)
	if err != nil {
		app.HandleError(w, "insert db failed", err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *Config) HandleError(w http.ResponseWriter, message string, err error) {
	app.Log.Error(err, "marshal request failed")
	app.errorJSON(w, err)
}
