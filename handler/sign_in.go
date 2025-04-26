package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/affordAbode_backend/model"
)

type signinReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Signin used to authenticate extant user
func (h *Handler) SignIn(c *gin.Context) {
	var req signinReq

	if ok := bindData(c, &req); !ok {
		fmt.Println("binding data unsuccessful")
		return
	}

	// inject user request to user to use for User service layer
	u := &model.Users{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	uNew, affordAbodeErr := h.UserService.Signin(ctx, u)

	if affordAbodeErr != nil {
		log.Printf("Failed to sign in user: %v\n", affordAbodeErr)
		// affordAbodeErr := errors.NewInternalServerError("Sign in error")
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	tokens, affordAbodeErr := h.TokenService.NewPairForUser(ctx, uNew, "")
	if affordAbodeErr != nil {
		log.Printf("Failed to create tokens for user: %v\n", affordAbodeErr)
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
