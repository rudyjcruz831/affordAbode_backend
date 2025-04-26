package services

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/google/uuid"
	"github.com/rudyjcruz831/affordAbode_backend/model"
	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
)

// TokenService used for injecting an implementation of TokenRepository
// for use in service methods along with keys and secrets for
// signing JWTs
type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationsSecs     int64
	RefreshExpirationSecs int64
}

// TSConfig will hold repositories that will eventually be injected into this
// this service layer
type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationsSecs     int64
	RefreshExpirationSecs int64
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationsSecs:     c.IDExpirationsSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewPairForUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (s *tokenService) NewPairForUser(ctx context.Context, u *model.Users, prevTokenID string) (*model.TokenPair, *errors.AffordAbodeError) {

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if affordAbodeErr := s.TokenRepository.DeleteRefreshToken(ctx, u.ID, prevTokenID); affordAbodeErr != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.ID, prevTokenID)
			return nil, affordAbodeErr
		}
	}

	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationsSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.ID, err.Error())
		affordAbodeErr := errors.NewInternalServerError("error generating id token")
		return nil, affordAbodeErr
	}

	refreshToken, err := generateRefreshToken(u.ID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.ID, err.Error())
		affordAbodeErr := errors.NewInternalServerError("error creating refresh Token")
		return nil, affordAbodeErr
	}

	// store refresh tokens by calling TokenRepository methods
	// set freshly minted refresh token to valid list
	if affordAbodeErr := s.TokenRepository.SetRefreshToken(ctx, u.ID, refreshToken.ID.String(), refreshToken.ExpiresIn); affordAbodeErr != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.ID, affordAbodeErr)
		return nil, errors.NewInternalServerError("")
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SS: idToken.SS, ID: idToken.ID},
		RefreshToken: model.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID, UID: u.ID},
	}, nil
}

// Signout reaches out to the repository layer to delete all valid tokens for a user
// TODO: test if correctly return id of idToken to use in blacklist add testing to token layer for BlackedListed
func (s *tokenService) Signout(ctx context.Context, uid string) *errors.AffordAbodeError {
	if affordAbodeErr := s.TokenRepository.DeleteUserRefreshTokens(ctx, uid); affordAbodeErr != nil {
		return affordAbodeErr
	}
	return nil
}

// ValidateIDToken validates the id token jwt string
// It returns the user extract from the IDTokenCustomClaims
func (s *tokenService) ValidateIDToken(tokenString string) (*model.Users, string, *errors.AffordAbodeError) {
	claims, err := validateIDToken(tokenString, s.PubKey) // uses public RSA key

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, "", errors.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, claims.Id, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, *errors.AffordAbodeError) {
	// validate actual JWT with string a secret
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, errors.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, errors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &model.RefreshToken{
		SS:  tokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}
