package tweetgo

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/schema"
)

// Client is the Twitter API client which will make signed requests to twitter
type Client struct {
	OAuthConsumerKey       string
	OAuthConsumerSecret    string
	OAuthAccessToken       string
	OAuthAccessTokenSecret string
	HTTPClient             requestMaker
	Noncer                 nonceMaker
	Timer                  currentTimer
}

type noncer struct{}

func (n noncer) Generate() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 48)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}

	return string(b)
}

type timer struct{}

func (t timer) GetCurrentTime() int64 {
	return time.Now().Unix()
}

func NewClient(oauthConsumerKey, oauthConsumerSecret string) Client {
	return Client{
		OAuthConsumerKey:    oauthConsumerKey,
		OAuthConsumerSecret: oauthConsumerSecret,
		HTTPClient:          &http.Client{},
		Noncer:              noncer{},
		Timer:               timer{},
	}
}

func (c *Client) SetAccessKeys(oauthAccessToken, oauthAccessTokenSecret string) {
	c.OAuthAccessToken = oauthAccessToken
	c.OAuthAccessTokenSecret = oauthAccessTokenSecret
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

// StatusesUserTimelineGet will get a users timeline and return an array of tweets
func (c Client) StatusesUserTimelineGet(input StatusesUserTimelineInput) ([]StatusesUserTimelineOutput, error) {
	uri := "https://api.twitter.com/1.1/statuses/user_timeline.json"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodGet, uri, params)
	if err != nil {
		return []StatusesUserTimelineOutput{}, err
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []StatusesUserTimelineOutput{}, err
	}

	var output []StatusesUserTimelineOutput
	json.Unmarshal(resBytes, &output)

	return output, nil
}
