package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tokenReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *Handler) Tokens(c *gin.Context) {
	fmt.Print("Tokens handler\n")
	var req tokenReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	//verify refresh JWT
	refreshToken, affordAbodeErr := h.TokenService.ValidateRefreshToken(req.RefreshToken)
	if affordAbodeErr != nil {
		log.Printf("Failed to validate refresh token: %v\n", affordAbodeErr)
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	// get up-to-date user
	u, affordAbodeErr := h.UserService.Get(ctx, refreshToken.UID)

	if affordAbodeErr != nil {
		log.Printf("Failed to get user: %v\n", affordAbodeErr)
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	tokens, affordAbodeErr := h.TokenService.NewPairForUser(ctx, u, refreshToken.ID.String())

	if affordAbodeErr != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", u, affordAbodeErr)
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokens,
	})
}
