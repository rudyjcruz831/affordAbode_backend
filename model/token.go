package model

import "github.com/google/uuid"

// RefreshToken stores token properties that
// are accessed in multiple application layers
type RefreshToken struct {
	SS  string    `json:"refreshToken"`
	ID  uuid.UUID `json:"-"`
	UID string    `json:"-"`
}

// IDToken stores token properties that
// are accessed in multiple application layers
type IDToken struct {
	SS string    `json:"idToken"`
	ID uuid.UUID `json:"-"`
}

// TokenPair used for returning pairs of id and refresh tokens
type TokenPair struct {
	IDToken
	RefreshToken
}
