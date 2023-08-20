package config

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
