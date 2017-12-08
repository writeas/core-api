package as

import (
	"net/http"
	"time"
)

// defaultHTTPTimeout is the default http.Client timeout.
const defaultHTTPTimeout = 10 * time.Second

// ClientConfig contains API configuration
type ClientConfig struct {
	// BaseURL of the API we're talking to
	BaseURL string

	// Client making requests to the API
	Client *http.Client
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		BaseURL: "",
		Client:  &http.Client{Timeout: defaultHTTPTimeout},
	}
}
