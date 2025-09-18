package commonmodel

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error" example:"error message"`
	Code    int    `json:"code,omitempty" example:"400"`
	Details string `json:"details,omitempty" example:"Invalid user id"`
}
