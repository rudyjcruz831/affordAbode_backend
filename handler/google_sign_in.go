package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthCode ...
type AuthCodeReq struct {
	Code string `json:"code" binding:"required"`
}

func (h *Handler) GoogleSignin(c *gin.Context) {
	// authCode struct
	var req AuthCodeReq

	if ok := bindData(c, &req); !ok {
		log.Println("binding data unsuccess")
		return
	}
	ctx := c.Request.Context()
	u, mathSheetsErr := h.UserService.GoogleSignin(ctx, req.Code)
	if mathSheetsErr != nil {
		log.Printf("Failed to sign up user using google: %v\n", mathSheetsErr)
		c.JSON(mathSheetsErr.Status, mathSheetsErr)
		return
	}

	fmt.Println("User signed in successfully:", u)

	// create token pair as strings
	tokens, mathSheetsErr := h.TokenService.NewPairForUser(ctx, u, "")
	if mathSheetsErr != nil {
		log.Printf("Failed to create tokens for user: %v\n", mathSheetsErr)

		// logic to go into database and delete user
		// when token NewPairForUser faileds
		err := h.UserService.DeleteUser(ctx, u.ID)
		if err != nil {
			log.Printf("Token Creation failed and Deleting user also failed: %v", err)
		}

		c.JSON(mathSheetsErr.Status, mathSheetsErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
