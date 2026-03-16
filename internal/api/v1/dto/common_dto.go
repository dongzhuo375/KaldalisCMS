package dto

// ErrorResponse is a uniform API error envelope.
type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

// MessageResponse is a generic API success envelope.
type MessageResponse struct {
	Message string `json:"message"`
}
