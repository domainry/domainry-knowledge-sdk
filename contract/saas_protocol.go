package contract

import (
	"encoding/json"
)

const (
	SaaSDiscoveryPath = "/knowledge/v1/discovery"
	SaaSInvokePath    = "/knowledge/v1/invoke"
)

// SaaSRequest is the bounded transport envelope used by the official remote
// client. Operation names are owned by this SDK; Input never contains a Store,
// database handle, SQL transaction or deployment credential.
type SaaSRequest struct {
	Operation string          `json:"operation"`
	Input     json.RawMessage `json:"input,omitempty"`
	Grants    []SaaSGrant     `json:"grants,omitempty"`
}

// SaaSResponse keeps implementation errors and response bodies behind a
// stable SDK-owned boundary. Implementations must not place provider, SQL or
// credential details in Message.
type SaaSResponse struct {
	Result    json.RawMessage `json:"result,omitempty"`
	Error     *SaaSError      `json:"error,omitempty"`
	Challenge *SaaSChallenge  `json:"challenge,omitempty"`
}

type SaaSError struct {
	Class     string `json:"class"`
	Code      string `json:"code"`
	Message   string `json:"message,omitempty"`
	Retryable bool   `json:"retryable,omitempty"`
}

type SaaSChallenge struct {
	Token string          `json:"token"`
	Kind  string          `json:"kind"`
	Input json.RawMessage `json:"input"`
}

type SaaSGrant struct {
	Token  string          `json:"token"`
	Result json.RawMessage `json:"result,omitempty"`
}
