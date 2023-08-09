package main

import (
	"net/http"
	"random-service/cmd/api/models"
	dbmodels "random-service/cmd/database/models"
)

func (app *Config) SaveLog(w http.ResponseWriter, r *http.Request) {
	var request models.JSONPayload
	err := app.readJSON(w, r, &request)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	entry := dbmodels.RandomData{
		Data: request.Data,
	}
	app.log.PrintLogs("Before save DB")

	err = app.db.Insert(entry)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	app.log.PrintLogs(request)

	resData := models.ResponseData{
		Data: "Called protocol is: JSON",
	}
	res := models.JSONResponse{
		Error:   false,
		Message: "Response from random-service",
		Data:    resData,
	}
	app.writeJSON(w, http.StatusAccepted, res)
}
