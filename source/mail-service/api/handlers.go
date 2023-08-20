package api

import (
	"net/http"

	"mailer-service/internal/mailer"
	"mailer-service/pkg/utils"
)

func (app *Server) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := utils.ReadJSON(w, r, &requestPayload)
	if err != nil {
		utils.HandleError(w, app.Log, "failed to read mail request", err)
		return
	}

	msg := mailer.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		utils.HandleError(w, app.Log, "send smtp message failed", err)
		return
	}

	payload := utils.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	err = utils.WriteJSON(w, http.StatusAccepted, payload)
	if err != nil {
		app.Log.Error(err, "failed to write mail response")
	}
}
