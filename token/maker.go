package token

import "time"

// Maker interface to managing tokens
type Maker interface {
	// CreateToken creates token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if a token is valid or not
	VerifyToken(token string) (*Payload, error)
}
