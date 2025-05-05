package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rudyjcruz831/affordAbode_backend/model"
)

type ReqRentalListings struct {
	Address   string  `json:"address"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	RentPrice float64 `json:"rentPrice"`
	Bedrooms  int     `json:"bedrooms"`
	Bathrooms int     `json:"bathrooms"`
}

type RentalResponse struct {
	Listings []ReqRentalListings `json:"listings"`
}

func fetchRentalListings(l *model.Listing) ([]model.Listing, error) {
	url := fmt.Sprintf("https://api.rentcast.io/v1/listings?city=%s&state=%s", l.City, l.State)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 	curl --request GET \
	//   --url 'https://api.rentcast.io/v1/properties?city=Austin&state=TX&limit=20' \
	//   --header 'Accept: application/json' \
	//   --header 'X-Api-Key: YOUR_API_KEY'
	// get apiKey form env variable
	apiKey := os.Getenv("RENTCAST_API_KEY")
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("API request failed with status:", resp.Status)
		// Read the response body to get more details about the error
		log.Println("Response Body:", resp.Body)
		// You can also log the response body if needed
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Response Body:", string(body))
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	fmt.Printf("Response Status: %v\n", resp)
	// var rentalData RentalResponse
	// if err := json.NewDecoder(resp.Body).Decode(&rentalData); err != nil {
	// 	return nil, err
	// }

	listings := make([]model.Listing, 0)
	// listings := make([]model.Listing, 0, len(rentalData.Listings))
	// for _, rental := range rentalData.Listings {
	// 	listings = append(listings, model.Listing{
	// 		Address:   rental.Address,
	// 		City:      rental.City,
	// 		State:     rental.State,
	// 		RentPrice: rental.RentPrice,
	// 		Bedrooms:  rental.Bedrooms,
	// 		Bathrooms: rental.Bathrooms,
	// 	})
	// }

	return listings, nil
}
