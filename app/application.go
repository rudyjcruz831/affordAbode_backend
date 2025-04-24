package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartApp() {
	// Create a new instance of the app
	// fmt.Println("Hello, World!")

	//injection other services and add any env variables TODO: may need this for injecting the serverses later
	ds, err := initDS()
	if err != nil {
		log.Fatalf("Failed to initialize data source: %v\n", err)
	}
	router, err := inject(ds) // remove ds for now becuse we don't have no db for now [router, err := inject(ds)]
	if err != nil {
		log.Fatalf("Failure to inject data source: %v\n", err)
	}

	//grabbing port from env for running server local or other host
	port := os.Getenv("USER_API_PORT")
	//if port env is empty the make it default 50052
	if port == "" {
		port = "50052"
	}

	// Graceful server shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		// running the server on localhost with given port
		//------------------------- This CODE HERE is for TLS-------------------------
		// if err := srv.ListenAndServeTLS("./server.crt", "./server.pem"); err != nil && err != http.ErrServerClosed{
		// 	log.Fatalf("Failed to intialized server: %v\n", err)
		// }

		// This code is without TLS
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to intialized server: %v\n", err)
		}
	}()

	fmt.Printf("Listining on port %s\n", port)

	//wait for kill signal of channel
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This block until a signal is passed into the quit channel
	<-quit
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

}
