package models

// ErrorResponse - is the basic error response structure
type ErrorResponse struct {
	Error ErrorResponseFormat `json:"error"`
}

// ErrorResponseFormat - is a simple error format according to appventurez standards...
type ErrorResponseFormat struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	MessageCode string `json:"messageCode"`
}
