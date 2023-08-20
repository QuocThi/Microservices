package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"authentication/internal/config"
	"authentication/pkg/utils"
)

func (app *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
	request := config.RequestPayload{}
	err := utils.ReadJSON(w, r, &request)
	if err != nil {
		app.Log.Error(err, "parse authenticate request failed")
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(request.Email)
	if err != nil {
		app.Log.Error(err, "invalid credentials")
		utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(request.Password)
	if err != nil || !valid {
		app.Log.Error(err, "wrong password")
		utils.ErrorJSON(w, errors.New("wrong password"), http.StatusBadRequest)
		return
	}

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.Log.Error(err, "failed to log authenticate request")
		utils.ErrorJSON(w, err)
		return
	}

	payload := utils.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	utils.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Server) logRequest(name, data string) error {
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
