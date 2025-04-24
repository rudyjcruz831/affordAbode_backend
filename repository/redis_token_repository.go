package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rudyjcruz831/affordAbode_backend/model"
	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
)

// redisTokenRepository is data/repository implementation
// of service layer TokenRepository
type redisTokenRepository struct {
	Redis *redis.Client
}

// NewTokenRepository is a factory for initializing User Repositories
func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken stores a refresh token with an expiry time
func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.AffordAbodeError {
	// panic("Set Refresh actuall iplmentation to interact with redis")
	// We'll store userID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return errors.NewInternalServerError("")
	}
	return nil
}

// DeleteRefreshToken used to delete old  refresh tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) *errors.AffordAbodeError {
	// panic("jlkjflkajsf")
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := r.Redis.Del(ctx, key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return errors.NewInternalServerError("")
	}

	// // Val returns count of deleted keys
	// // If no value was deleted, the refresh token is invalid
	// fmt.Print(result.Val(), "\n")
	// if result.Val() < 1 {
	// 	log.Printf("Refresh token to redis for userID: %s:%s does not exist\n", userID, tokenID)
	// 	return errors.NewAuthorization("Invalid refresh token")
	// }

	return nil
}

// DeleteUserRefreshTokens looks for all tokens beginning with
// userID and scans to delete them in a non-blocking fashion
func (r *redisTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) *errors.AffordAbodeError {
	pattern := fmt.Sprintf("%s*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	for iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("Failed to delete refresh token: %s\n", iter.Val())
			failCount++
		}
	}

	// check last value
	if err := iter.Err(); err != nil {
		log.Printf("Failed to delete refresh token token iter Error: %s\n", iter.Val())
	}

	if failCount > 0 {
		return errors.NewInternalServerError("")
	}

	return nil
}

// TokenBlackedListed add tokens to redis to when signed out
// func (r *redisTokenRepository) TokenBlackedListed(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) *errors.AffordAbodeError {
// 	// panic("")
// 	// We'll store userID with token id so we can scan (non-blocking)
// 	// over the user's tokens and delete them in case of token leakage
// 	key := fmt.Sprintf("%s:%s", userID, tokenID)
// 	if err := r.Redis.Set(ctx, key, 1, expiresIn).Err(); err != nil {
// 		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
// 		return errors.NewInternalServerError("")
// 	}
// 	return nil
// }
