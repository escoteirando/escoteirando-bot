package domain

import "time"

type CachedAuthKey struct {
	ValidUntil time.Time
	AuthKey    string
}
