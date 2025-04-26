package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/rudyjcruz831/affordAbode_backend/model"
	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s model.TokenService, u model.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		h := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := errors.NewBadRequestError("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status, gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}

			// otherwise error type is unknown
			err := errors.NewInternalServerError("trying to bind")
			c.JSON(err.Status, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := errors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")
			c.JSON(err.Status, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		// validate ID token here
		// Validate the ID token and retrieve the user and token ID.
		// Parameters:
		// - idTokenHeader: The ID token extracted from the request header.
		// Returns:
		// - user: The validated user.
		// - tokenId: The ID of the token.
		// - err: Any error that occurred during validation.
		user, _, err := s.ValidateIDToken(idTokenHeader[1])
		if err != nil {
			err := errors.NewAuthorization("Provided token is invalid")
			c.JSON(err.Status, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		// checking if the user with that ID exist if he does not return error
		if _, affrodAbodeErr := u.Get(ctx, user.ID); affrodAbodeErr != nil {
			log.Printf("error: %v\n", affrodAbodeErr)
			c.JSON(affrodAbodeErr.Status, affrodAbodeErr)
			c.Abort()
			return
		}

		c.Set("user", user)
		// c.Set("token", idTokenHeader[1])

		// Continue to next handler
		c.Next()
	}
}
