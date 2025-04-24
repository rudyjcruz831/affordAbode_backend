package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rudyjcruz831/affordAbode_backend/handler"
	"github.com/rudyjcruz831/affordAbode_backend/repository"
	"github.com/rudyjcruz831/affordAbode_backend/services"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data source...")

	/*
		Repository Layer
	*/
	// TODO : add this when creating the AWS s3 bucket to store worksheets created history
	// workSheetBucketName := os.Getenv("AWS_FILE_BUCKET")
	userRepository := repository.NewUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(d.RedisClient)

	// imageRepository := repository.NewWorkSheetRepository(d.StorageClient, workSheetBucketName)

	/*
		Service Layer
	*/

	userServcie := services.NewUserService(&services.USConfig{
		UserRepository: userRepository,
	})

	// load rsa keys
	// TODO : what do i need this keys for?
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// load expiration lengths from env variables and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenService := services.NewTokenService(&services.TSConfig{
		TokenRepository:       tokenRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationsSecs:     idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// initialize gin.Engine
	router := gin.Default()

	// read in AffordAbode API url
	baseURL := os.Getenv("AFFORDABODE_API_URL")

	// handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	// ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	// }

	handler.NewHandler(&handler.Config{
		R:            router,
		BaseURL:      baseURL,
		UserService:  userServcie,
		TokenService: tokenService,
		// TimeoutDurations: time.Duration(time.Duration(ht) * time.Second),
		MaxBodyBytes: 1024 * 1024 * 1024,
	})

	return router, nil
}
