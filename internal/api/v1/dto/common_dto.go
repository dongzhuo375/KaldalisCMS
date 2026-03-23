package dto

// ErrorResponse is a uniform API error envelope.
type ErrorResponse struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

// MessageResponse is a generic API success envelope.
type MessageResponse struct {
	Message string `json:"message"`
}

// HealthCheckResult describes one dependency check in readiness probes.
type HealthCheckResult struct {
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// HealthResponse is a stable and extensible probe response contract.
type HealthResponse struct {
	Status string                       `json:"status"`
	Mode   string                       `json:"mode"`
	Checks map[string]HealthCheckResult `json:"checks"`
}
