package twitter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/schema"
)

// Client is the Twitter API client which will make signed requests to twitter
type Client struct {
	OAuthConsumerKey       string
	OAuthConsumerSecret    string
	OAuthAccessToken       string
	OAuthAccessTokenSecret string
	HTTPClient             http.Client
}

// OAuthRequestTokenGet will return an oauth_token and oauth_token_secret
func (c Client) OAuthRequestTokenGet(input OAuthRequestTokenInput) (OAuthRequestTokenOutput, error) {
	uri := "https://api.twitter.com/oauth/request_token"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodPost, uri, params)
	if err != nil {
		return OAuthRequestTokenOutput{}, err
	}
	defer res.Body.Close()

	values, err := bodyToValues(res.Body)
	if err != nil {
		return OAuthRequestTokenOutput{}, err
	}

	var output = OAuthRequestTokenOutput{}
	var decoder = schema.NewDecoder()
	err = decoder.Decode(&output, values)

	return output, nil
}

// OAuthAccessTokenGet will exchange a temporary access token for a permanent one
func (c Client) OAuthAccessTokenGet(input OAuthAccessTokenInput) (OAuthAccessTokenOutput, error) {
	uri := "https://api.twitter.com/oauth/access_token"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodPost, uri, params)
	if err != nil {
		return OAuthAccessTokenOutput{}, err
	}
	defer res.Body.Close()

	values, err := bodyToValues(res.Body)
	if err != nil {
		return OAuthAccessTokenOutput{}, err
	}

	var output = OAuthAccessTokenOutput{}
	var decoder = schema.NewDecoder()
	err = decoder.Decode(&output, values)

	return output, nil
}

// StatusesUpdatePost will post a status update to twitter
func (c Client) StatusesUpdatePost(input StatusesUpdateInput) (StatusesUpdateOutput, error) {
	uri := "https://api.twitter.com/1.1/statuses/update.json"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodPost, uri, params)
	if err != nil {
		return StatusesUpdateOutput{}, err
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return StatusesUpdateOutput{}, err
	}

	output := StatusesUpdateOutput{}
	json.Unmarshal(resBytes, &output)

	return output, nil
}

// StatusesFilterPostRaw will get a streaming list of tweets and return the raw http response for streaming
func (c Client) StatusesFilterPostRaw(input StatusesFilterInput) (*http.Response, error) {
	uri := "https://stream.twitter.com/1.1/statuses/filter.json"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodPost, uri, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}
