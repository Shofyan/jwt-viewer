package models

import "time"

// DecodeRequest represents a request to decode a JWT
type DecodeRequest struct {
	Token string `json:"token" binding:"required"`
}

// DecodeResponse represents the decoded JWT parts
type DecodeResponse struct {
	Header    map[string]interface{} `json:"header"`
	Payload   map[string]interface{} `json:"payload"`
	Signature string                 `json:"signature"`
	Error     string                 `json:"error,omitempty"`
}

// EncodeRequest represents a request to encode/generate a JWT
type EncodeRequest struct {
	Header    map[string]interface{} `json:"header" binding:"required"`
	Payload   map[string]interface{} `json:"payload" binding:"required"`
	Secret    string                 `json:"secret" binding:"required"`
	Algorithm string                 `json:"algorithm"` // HS256, HS384, HS512
}

// EncodeResponse represents the generated JWT
type EncodeResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// VerifyRequest represents a request to verify a JWT
type VerifyRequest struct {
	Token  string `json:"token" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

// VerifyResponse represents the verification result
type VerifyResponse struct {
	Valid   bool                   `json:"valid"`
	Message string                 `json:"message"`
	Claims  map[string]interface{} `json:"claims,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// ClaimInfo represents human-readable claim information
type ClaimInfo struct {
	Exp       *time.Time `json:"exp,omitempty"`
	ExpString string     `json:"exp_string,omitempty"`
	Iat       *time.Time `json:"iat,omitempty"`
	IatString string     `json:"iat_string,omitempty"`
	Nbf       *time.Time `json:"nbf,omitempty"`
	NbfString string     `json:"nbf_string,omitempty"`
	Iss       string     `json:"iss,omitempty"`
	Aud       string     `json:"aud,omitempty"`
	Sub       string     `json:"sub,omitempty"`
	IsExpired bool       `json:"is_expired"`
}
