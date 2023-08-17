package models

type JSONPayload struct {
	Data string `json:"data"`
}

type JSONResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type ResponseData struct {
	Data string `json:"data"`
}
