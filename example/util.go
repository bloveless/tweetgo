package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"tweetgo"
)

type config struct {
	OAuthConsumerKey       string `json:"oauth_consumer_key"`
	OAuthConsumerSecret    string `json:"oauth_consumer_secret"`
	OAuthAccessToken       string `json:"oauth_access_token"`
	OAuthAccessTokenSecret string `json:"oauth_access_token_secret"`
}

func loadConfig() (config, error) {
	wd, _ := os.Getwd()

	fmt.Printf("Loading config from %s/config.json\n", wd)
	configBytes, err := ioutil.ReadFile(wd+"/config.json")
	if err != nil {
		return config{}, err
	}

	c := config{}
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		return config{}, err
	}

	fmt.Printf("Config loaded successfully: %+v\n", c)
	return c, nil
}

func saveConfig(c config) error {
	wd, _ := os.Getwd()

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	fmt.Printf("Saving new config to %s/config.json\n", wd)
	f, err := os.OpenFile(wd+"/config.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully wrote new config")
	return nil
}

func getTwitterClient(c config) tweetgo.Client {
	tc := tweetgo.NewClient(
		c.OAuthConsumerKey,
		c.OAuthConsumerSecret,
	)

	tc.SetAccessKeys(c.OAuthAccessToken, c.OAuthAccessTokenSecret)

	return tc
}
