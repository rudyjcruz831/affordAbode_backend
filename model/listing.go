package model

type Listing struct {
	Address   string  `json:"address" binding:"required"`
	City      string  `json:"city" binding:"required"`
	State     string  `json:"state" binding:"required"`
	RentPrice float64 `json:"rent_price" binding:"required"`
	Bedrooms  int     `json:"bedrooms" binding:"required"`
	Bathrooms int     `json:"bathrooms" binding:"required"`
}
