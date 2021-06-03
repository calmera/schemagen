package schemagen

type Schema struct {
	Schema  []byte `json:"schema"`
	Subject string `json:"subject"`
	Version int    `json:"version"`
	ID      int    `json:"id,omitempty"`
}
