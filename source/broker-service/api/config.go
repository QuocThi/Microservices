package api

// RequestPayload describes the JSON that this service accepts as an HTTP Post request
type RequestPayload struct {
	Action string        `json:"action"`
	Auth   AuthPayload   `json:"auth,omitempty"`
	Log    LogPayload    `json:"log,omitempty"`
	Mail   MailPayload   `json:"mail,omitempty"`
	Random RandomPayload `json:"random,omitempty"`
}

type RandomPayload struct {
	Data string `json:"data"`
}

// MailPayload is the embedded type (in RequestPayload) that describes an email message to be sent
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// AuthPayload is the embedded type (in RequestPayload) that describes an authentication request
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LogPayload is the embedded type (in RequestPayload) that describes a request to log something
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
