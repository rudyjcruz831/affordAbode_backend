package services

import (
	"context"
	"log"

	"github.com/rudyjcruz831/affordAbode_backend/model"
	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
)

// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository model.UserRepository
	// DocsRepository  model.DocsRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository model.UserRepository
	// DocsRepository  model.DocsRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
		// DocsRepository:  c.DocsRepository,
	}
}

// Get retrieves a user based on their ID
func (s *userService) Get(ctx context.Context, id string) (*model.Users, *errors.AffordAbodeError) {
	u, err := s.UserRepository.FindByID(ctx, id)
	return u, err
}

// Signup reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signup(ctx context.Context, u *model.Users) *errors.AffordAbodeError {
	pw, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return errors.NewInternalServerError("")
	}
	u.Password = pw

	if mathShtErr := s.UserRepository.Create(ctx, u); mathShtErr != nil {
		log.Printf("UserRepository return error: %v", mathShtErr)
		return mathShtErr
	}

	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishUserUpdated(u, true)

	// if err != nil {
	// 	return nil, apperrors.NewInternal()
	// }

	return nil
}

// Signin reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signin(ctx context.Context, u *model.Users) (*model.Users, *errors.AffordAbodeError) {
	// panic("Sign In Method not implemented")
	uFetched, MathShtErr := s.UserRepository.FindByEmail(ctx, u.Email)

	// Will return NotAuthorized to client to omit details of why
	if MathShtErr != nil {
		log.Printf("FindByEmail return error : %v", MathShtErr)
		return nil, errors.NewAuthorization("Invalid email and password combination")
	}

	match, err := comparePasswords(uFetched.Password, u.Password)

	if err != nil {
		log.Printf("comparePassword return error %v", err)
		return nil, errors.NewInternalServerError("")
	}

	if !match {
		log.Println("Match was false return error")
		return nil, errors.NewAuthorization("Invalid email and password combination")
	}

	return uFetched, nil
}

// Update Details reaches out
func (s *userService) UpdateDetails(ctx context.Context, u *model.Users) *errors.AffordAbodeError {
	// Update user in UserRepository
	MathShtErr := s.UserRepository.Update(ctx, u)

	if MathShtErr != nil {
		return MathShtErr
	}

	// // Publish user updated nats streaming server // kafca
	// err = s.EventsBroker.PublishUserUpdated(u, false)
	// if err != nil {
	// 	return errors.NewInternal()
	// }

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) *errors.AffordAbodeError {
	// panic("Delete user service")
	mathSheetsErr := s.UserRepository.Delete(ctx, id)
	return mathSheetsErr
}

// Using google to sign in user if no user will be created calling the repo
func (s *userService) GoogleSignin(ctx context.Context, code string) (*model.Users, *errors.AffordAbodeError) {
	// panic("Google Sing In")
	// TODO - testing for this service
	_, u, afordAbodeErr := auth(code)
	if afordAbodeErr != nil {
		// c.JSON(afordAbodeErr.Status, afordAbodeErr)
		return nil, afordAbodeErr
	}

	uFetched, afordAbodeErr := s.UserRepository.FindByEmail(ctx, u.Email)
	if afordAbodeErr != nil {
		if afordAbodeErr.Status == 404 {
			u.Role = "user"
			if afordAbodeErr = s.UserRepository.Create(ctx, u); afordAbodeErr != nil {
				return nil, afordAbodeErr
			}
			return u, nil
		} else {
			return nil, afordAbodeErr
		}
	}

	return uFetched, nil
}

func (s *userService) Listing(ctx context.Context, listing *model.Listing) ([]model.Listing, *errors.AffordAbodeError) {
	// panic("Listing service")
	// TODO - testing for this service
	list, err := fetchRentalListings(listing)
	if err != nil {
		log.Printf("Error fetching rental listing: %v\n", err)
		return nil, errors.NewInternalServerError("Error fetching rental listing")
	}

	return list, nil

}
