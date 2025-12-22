package heysender

import (
	"fmt"
	"net/http"
)

type HeysenderClient struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	HTTPClient *http.Client
}

// HeysenderError represents an API error
type HeysenderError struct {
	Message    string
	StatusCode int
}

func (e *HeysenderError) Error() string {
	return fmt.Sprintf("Heysender API Error (%d): %s", e.StatusCode, e.Message)
}
