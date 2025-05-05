package services

import (
	"context"
	"encoding/json"
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

	return &response, &u, nil
}

func getConfig() *oauth2.Config {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	return conf
}

func getTokenFromAuthCode(authCode string) (*oauth2.Token, *errors.AffordAbodeError) {
	conf := getConfig()

	// Exchange authorization code for token
	tok, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
		affordAbodeErr := errors.UnauthorizedError("getting token from authCode: " + err.Error())
		return nil, affordAbodeErr
	}

	return tok, nil
}

func getUserInfoFromToken(token *oauth2.Token) ([]byte, *errors.AffordAbodeError) {
	conf := getConfig()

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		affordAbodeErr := errors.UnauthorizedError("getting user info from token: " + err.Error())
		return nil, affordAbodeErr
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		affordAbodeErr := errors.NewInternalServerError("error reading response body")
		return nil, affordAbodeErr
	}

	return data, nil
}
