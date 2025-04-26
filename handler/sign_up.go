package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/affordAbode_backend/model"
)

// signupReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
type signupReq struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signup handler
func (h *Handler) SignUp(c *gin.Context) {
	// panic("Sign up handler")
	// define a variable to which we'll bind incoming
	// json body, {email, password}
	var req signupReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	// inject user request to user to use for User service layer
	u := &model.Users{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}

	ctx := c.Request.Context()

	// UserService layer
	affordAbodeErr := h.UserService.Signup(ctx, u)
	if affordAbodeErr != nil {
		log.Printf("Failed to sign up user: %v\n", affordAbodeErr)
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	// create token pair as strings
	tokens, affordAbodeErr := h.TokenService.NewPairForUser(ctx, u, "")
	if affordAbodeErr != nil {
		log.Printf("Failed to create tokens for user: %v\n", affordAbodeErr)

		// logic to go into database and delete user
		// when token NewPairForUser faileds
		ctx = c.Request.Context()
		err := h.UserService.DeleteUser(ctx, u.ID)
		if err != nil {
			log.Printf("Token Creation failed and Deleting user also failed: %v", err)
		}

		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
