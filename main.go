package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	

	"github.com/michaeld3289/microservice_tutorial/coffee_ecom/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	products := handlers.NewProducts(logger)

	sm := http.NewServeMux()
	sm.Handle("/", products)

	// create a new server
	server := http.Server{
		Addr:         ":9090",      
		Handler:      sm,                // set the default handler
		ErrorLog:     logger,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		logger.Println("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}