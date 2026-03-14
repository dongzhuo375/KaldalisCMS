package dto

// ErrorResponse is a uniform API error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse is a generic API success envelope.
type MessageResponse struct {
	Message string `json:"message"`
}
