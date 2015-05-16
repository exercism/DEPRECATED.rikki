package main

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
)

// Auth is a simple shared secret.
// If the rikki project ever becomes a larger part of the
// exercism ecosystem, we should consider improving the security
// around it.
type Auth struct {
	Secret string
}

// NewAuth configures a shared secret.
// It looks in an environment variable for a pass phrase, and
// falls back to a common value for development if not set.
func NewAuth() *Auth {
	a := Auth{}
	a.Secret = os.Getenv("RIKKI_SECRET")
	if a.Secret == "" {
		// Use a default value in development mode.
		a.Secret = "I wish a robot would get elected president. That way, when he came to town, we could all take a shot at him and not feel too bad."
	}
	return &a
}

// Key returns the SHA1 hexdigest of the shared secret
// kind of stupid, since there's no salt
func (a *Auth) Key() string {
	hasher := sha1.New()
	hasher.Write([]byte(a.Secret))
	return hex.EncodeToString(hasher.Sum(nil))
}
