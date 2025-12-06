package handler

// RequestBody represents the JSON request body
type RequestBody struct {
	URL string `json:"url"`
}

// SuccessResponse represents a successful request creation response
type SuccessResponse struct {
	DocumentID string `json:"documentId"`
	Status     string `json:"status"`
	URL        string `json:"url"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
