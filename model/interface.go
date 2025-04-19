package model

import (
	"context"

	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
)

// UserService defines methods the handler layer expects
type UserService interface {
	// Get(ctx context.Context, id string) (*Users, *errors.MathSheetsError)
	// Signup(ctx context.Context, u *Users) *errors.MathSheetsError
	// Signin(ctx context.Context, u *Users) (*Users, *errors.MathSheetsError)
	// UpdateDetails(ctx context.Context, u *Users) *errors.MathSheetsError
	// DeleteUser(ctx context.Context, id string) *errors.MathSheetsError
	GoogleSignin(ctx context.Context, code string) (*Users, *errors.AffordAbodeError)
}
