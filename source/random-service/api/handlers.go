package api

import (
	"net/http"

	"random-service/api/models"
	dbmodels "random-service/internal/models"
	"random-service/pkg/utils"
)

func (app *Server) SaveRandom(w http.ResponseWriter, r *http.Request) {
	var request models.JSONPayload
	err := utils.ReadJSON(w, r, &request)
	if err != nil {
		utils.HandleError(w, app.Log, "parse request failed", err)
		return
	}

	entry := dbmodels.RandomData{
		Data: request.Data,
	}

	err = app.DB.Insert(entry)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resData := models.ResponseData{
		Data: "Called protocol is: JSON",
	}
	res := models.JSONResponse{
		Error:   false,
		Message: "Response from random-service",
		Data:    resData,
	}
	err = utils.WriteJSON(w, http.StatusAccepted, res)
	if err != nil {
		app.Log.Error(err, "failed to save random info")
	}
}
