package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tweetgo"
)

func auth() config {
	c, err := loadConfig()
	if err != nil {
		panic(err)
	}

	if c.OAuthAccessToken == "" || c.OAuthAccessTokenSecret == "" {
		fmt.Println("Unable to find access keys. Getting new ones")
		c = requestToken(c)
	} else {
		fmt.Println("Access keys already present")
	}

	return c
}

func requestToken(c config) config {
	tc := getTwitterClient(c)

	input := tweetgo.OAuthRequestTokenInput{
		OAuthCallback: tweetgo.String("http://127.0.0.1:3000/oauth_response"),
	}

	output, err := tc.OAuthRequestTokenGet(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Please authenticate by visiting this link https://api.twitter.com/oauth/authorize?oauth_token=" + output.OAuthToken)

	type oAuthResponse struct {
		OAuthToken    string
		OAuthVerifier string
	}

	done := make(chan oAuthResponse, 1)
	shutdown := make(chan bool, 1)

	router := http.NewServeMux()
	router.HandleFunc("/oauth_response", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access has been successfully recorded. You can return to the comfort of your terminal"))

		or := oAuthResponse{
			OAuthToken:    r.URL.Query().Get("oauth_token"),
			OAuthVerifier: r.URL.Query().Get("oauth_verifier"),
		}

		done <- or
	})

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	listenAddr := ":3000"
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	authResponse := oAuthResponse{}
	go func() {
		authResponse = <-done
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(shutdown)
	}()

	logger.Println("Server is waiting for auth token at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-shutdown
	logger.Println("Server stopped")

	accessTokenInput := tweetgo.OAuthAccessTokenInput{
		OAuthToken:    tweetgo.String(authResponse.OAuthToken),
		OAuthVerifier: tweetgo.String(authResponse.OAuthVerifier),
	}
	accessTokenOutput, err := tc.OAuthAccessTokenGet(accessTokenInput)
	if err != nil {
		panic(err)
	}

	updatedConfig := config{
		OAuthConsumerKey:       c.OAuthConsumerKey,
		OAuthConsumerSecret:    c.OAuthConsumerSecret,
		OAuthAccessToken:       accessTokenOutput.OAuthToken,
		OAuthAccessTokenSecret: accessTokenOutput.OAuthTokenSecret,
	}

	err = saveConfig(updatedConfig)
	if err != nil {
		panic(err)
	}

	return updatedConfig
}
