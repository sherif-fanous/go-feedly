package feedly

import "fmt"

// APIError represents a Feedly API Error response.
// https://developer.feedly.com/cloud/#client-errors
type APIError struct {
	ErrorID      string `json:"errorId"`
	ErrorMessage string `json:"errorMessage"`
}

// Error returns the string representation of an APIError.
func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorID, e.ErrorMessage)
}

func (e APIError) isEmpty() bool {
	return e.ErrorID == "" && e.ErrorMessage == ""
}

// relevantError returns any non-nil HTTP related error (Creating the request, getting the response, or decoding the
// response) if any. Otherwise return the decoded apiError (nil in case of no error).
func relevantError(httpError error, apiError *APIError) error {
	if httpError != nil {
		return httpError
	}

	if !apiError.isEmpty() {
		return apiError
	}

	return nil
}
