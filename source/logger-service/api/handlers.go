package api

import (
	"net/http"

	"log-service/internal/database"
	"log-service/pkg/utils"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload JSONPayload
	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.HandleError(w, app.Log, "marshal request failed", err)
		return
	}

	// insert data
	event := db.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.DB.Insert(event)
	if err != nil {
		utils.HandleError(w, app.Log, "insert db failed", err)
		return
	}

	resp := utils.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	err = utils.WriteJSON(w, http.StatusAccepted, resp)
	if err != nil {
		app.Log.Error(err, "failed to write rpc response")
	}
}
