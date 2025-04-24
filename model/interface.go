package model

import (
	"context"
	"time"

	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
)

// UserService defines methods the handler layer expects
type UserService interface {
	Get(ctx context.Context, id string) (*Users, *errors.AffordAbodeError)
	Signup(ctx context.Context, u *Users) *errors.AffordAbodeError
	Signin(ctx context.Context, u *Users) (*Users, *errors.AffordAbodeError)
	UpdateDetails(ctx context.Context, u *Users) *errors.AffordAbodeError
	DeleteUser(ctx context.Context, id string) *errors.AffordAbodeError
	GoogleSignin(ctx context.Context, code string) (*Users, *errors.AffordAbodeError)
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*Users, *errors.AffordAbodeError)
	Create(ctx context.Context, u *Users) *errors.AffordAbodeError
	FindByEmail(ctx context.Context, email string) (*Users, *errors.AffordAbodeError)
	Update(ctx context.Context, u *Users) *errors.AffordAbodeError
	Delete(ctx context.Context, id string) *errors.AffordAbodeError
	// UpdateImage(ctx context.Context, u *Users, imageURL string) (*Users, *errors.AffordAbodeError)
}

// TokenService defines methods handler layer expects to interact
// with in regards to producing JWT as string
type TokenService interface {
	NewPairForUser(ctx context.Context, u *Users, prevTokenID string) (*TokenPair, *errors.AffordAbodeError)
	Signout(ctx context.Context, uid string) *errors.AffordAbodeError
	ValidateIDToken(tokenString string) (*Users, string, *errors.AffordAbodeError)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, *errors.AffordAbodeError)
	// IsBlackedListed(ctx context.Context, uid string, tokenid string) *errors.AffordAbodeError
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.AffordAbodeError
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) *errors.AffordAbodeError
	DeleteUserRefreshTokens(ctx context.Context, userID string) *errors.AffordAbodeError
	// TokenBlackedListed(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.AffordAbodeError
}
