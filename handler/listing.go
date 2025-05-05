package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/affordAbode_backend/model"
)

type listingReq struct {
	Address   string  `json:"address" binding:""`
	City      string  `json:"city" binding:"required"`
	State     string  `json:"state" binding:"required"`
	RentPrice float64 `json:"rent_price" binding:""`
	Bedrooms  int     `json:"bedrooms" binding:""`
	Bathrooms int     `json:"bathrooms" binding:""`
}

func (h *Handler) CreateListing(c *gin.Context) {
	var req listingReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	listing := &model.Listing{
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		RentPrice: req.RentPrice,
		Bedrooms:  req.Bedrooms,
		Bathrooms: req.Bathrooms,
	}

	list, err := h.UserService.Listing(ctx, listing)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	for i := range list {
		fmt.Println(list[i])
	}
	// list[i].CreatedAt = time.Time{

	c.JSON(http.StatusOK, list)
}
