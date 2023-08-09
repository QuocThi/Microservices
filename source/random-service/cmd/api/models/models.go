package models

type JSONPayload struct {
	Data string `json:"data"`
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ResponseData struct {
	Data string `json:"data"`
}
