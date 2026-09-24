package remote

import (
	"net/http"
	"time"
)

type Config struct {
	Endpoint           string
	ServiceAccessToken string
	HTTPClient         *http.Client
	RequestTimeout     time.Duration
	MaxRequestBytes    int64
	MaxResponseBytes   int64
}

func normalizeConfig(value Config) Config {
	if value.RequestTimeout <= 0 {
		value.RequestTimeout = 20 * time.Second
	}
	if value.MaxRequestBytes <= 0 {
		value.MaxRequestBytes = 32 << 20
	}
	if value.MaxResponseBytes <= 0 {
		value.MaxResponseBytes = 32 << 20
	}
	return value
}
