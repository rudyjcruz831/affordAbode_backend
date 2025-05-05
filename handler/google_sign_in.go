package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/affordAbode_backend/model"
)

// SupabaseAuth represents the authentication data from Supabase
type SupabaseAuth struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
	User         struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user" binding:"required"`
}

// GoogleSignin handles the Supabase Google OAuth sign-in process
func (h *Handler) GoogleSignin(c *gin.Context) {
	var req SupabaseAuth

	if ok := bindData(c, &req); !ok {
		log.Println("binding data unsuccessful")
		return
	}

	ctx := c.Request.Context()

	// Create or update user in your database
	u := &model.Users{
		ID:        req.User.ID,
		Email:     req.User.Email,
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
	}

	// Create token pair for the user
	tokens, affordAbodeErr := h.TokenService.NewPairForUser(ctx, u, "")
	if affordAbodeErr != nil {
		log.Printf("Failed to create tokens for user: %v\n", affordAbodeErr)
		c.JSON(affordAbodeErr.Status, affordAbodeErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
		"user": gin.H{
			"id":        u.ID,
			"email":     u.Email,
			"firstName": u.FirstName,
			"lastName":  u.LastName,
		},
	})
}
