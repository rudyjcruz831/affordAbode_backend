package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rudyjcruz831/affordAbode_backend/model"
	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func auth(code string) (*model.Response, *model.Users, *errors.AffordAbodeError) {
	tok, affordAbodeErr := getTokenFromAuthCode(code)
	if affordAbodeErr != nil {
		return nil, nil, affordAbodeErr
	}

	tokenID, affordAbodeErr := getUserInfoFromToken(tok)
	if affordAbodeErr != nil {
		return nil, nil, affordAbodeErr
	}

	var tokID model.TokenID
	if err := json.Unmarshal(tokenID, &tokID); err != nil {
		affordAbodeErr := errors.NewInternalServerError("not able to parse token")
		return nil, nil, affordAbodeErr
	}

	u := model.Users{
		// ID:        tokID.ID,
		Email:     tokID.Email,
		FirstName: tokID.FirstName,
		LastName:  tokID.LastName,
	}

	response := model.Response{
		AccessToken: tok.AccessToken,
		ID:          tokID.ID,
		FirstName:   tokID.FirstName,
		LastName:    tokID.LastName,
		Email:       tokID.Email,
		TokenType:   tok.TokenType,
	}

	// fmt.Println(u)
	// fmt.Println(response)
	fmt.Println(response.TokenType)

	return &response, &u, nil
}

func getConfig() *oauth2.Config {
	// Get the Client Id and Client secret stored in enviroment variables
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	// fmt.Println(clientID, clientSecret)

	// Build auth configuration instance
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"profile", "email", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
	}
	// Don't know why but we need this to make it work
	conf.RedirectURL = "postmessage"

	return conf
}

func getTokenFromAuthCode(authCode string) (*oauth2.Token, *errors.AffordAbodeError) {
	conf := getConfig()

	// Exchange consumable authorization code for refresh token
	tok, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		affordAbodeErr := errors.UnauthorizedError("getting token form authCode ---" + err.Error())
		return nil, affordAbodeErr
	}

	return tok, nil
}

func getUserInfoFromToken(token *oauth2.Token) ([]byte, *errors.AffordAbodeError) {
	conf := getConfig()

	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		affordAbodeErr := errors.UnauthorizedError("getting user from token ---" + err.Error())
		return nil, affordAbodeErr
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		affordAbodeErr := errors.NewInternalServerError("Something went wrong on our end")
		return nil, affordAbodeErr
	}

	return data, nil
}
