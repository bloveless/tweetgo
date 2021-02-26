package tweetgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

// OAuthRequestTokenPost will return an oauth_token and oauth_token_secret
// https://developer.twitter.com/en/docs/basics/authentication/api-reference/request_token
func (c Client) OAuthRequestTokenPost(input OAuthRequestTokenInput) (OAuthRequestTokenOutput, error) {
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
	if err != nil {
		return OAuthRequestTokenOutput{}, err
	}

	return output, nil
}

// OAuthAccessTokenPost will exchange a temporary access token for a permanent one
// https://developer.twitter.com/en/docs/basics/authentication/api-reference/access_token
func (c Client) OAuthAccessTokenPost(input OAuthAccessTokenInput) (OAuthAccessTokenOutput, error) {
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
	if err != nil {
		return OAuthAccessTokenOutput{}, err
	}

	return output, nil
}

// ListsListGet will return all lists the authenticating user or specified user subscribes to, including thier own.
// https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/get-lists-list
func (c Client) ListsListGet(input ListsListInput) ([]ListsListOutput, error) {
	uri := "https://api.twitter.com/1.1/lists/list.json"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodGet, uri, params)
	if err != nil {
		return []ListsListOutput{}, err
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []ListsListOutput{}, err
	}

	output := []ListsListOutput{}
	err = json.Unmarshal(resBytes, &output)
	if err != nil {
		return []ListsListOutput{}, err
	}

	return output, nil
}

// ListsMembersGet will return the members of a specified list
// https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/get-lists-members
func (c Client) ListsMembersGet(input ListsMembersInput) (ListsMembersOutput, error) {
	uri := "https://api.twitter.com/1.1/lists/members.json"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodGet, uri, params)
	if err != nil {
		return ListsMembersOutput{}, err
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ListsMembersOutput{}, err
	}

	output := ListsMembersOutput{}
	err = json.Unmarshal(resBytes, &output)
	if err != nil {
		return ListsMembersOutput{}, err
	}

	return output, nil
}

// ListsMembersShowGet will check if the specified users is a member of the specified list
// https://developer.twitter.com/en/docs/accounts-and-users/create-manage-lists/api-reference/get-lists-members-show
func (c Client) ListsMembersShowGet(input ListsMembersShowInput) (ListsMembersShowOutput, error) {
	uri := "https://api.twitter.com/1.1/lists/members/show.json"
	params := processParams(input)

	res, err := c.executeRequest(http.MethodGet, uri, params)
	if err != nil {
		return ListsMembersShowOutput{}, err
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ListsMembersShowOutput{}, err
	}

	output := ListsMembersShowOutput{}
	err = json.Unmarshal(resBytes, &output)
	if err != nil {
		return ListsMembersShowOutput{}, err
	}

	return output, nil
}

// StatusesUpdatePost will post a status update to twitter
// https://developer.twitter.com/en/docs/tweets/post-and-engage/api-reference/post-statuses-update
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
	err = json.Unmarshal(resBytes, &output)
	if err != nil {
		return StatusesUpdateOutput{}, err
	}

	return output, nil
}

// StatusesFilterPostRaw will get a streaming list of tweets and return the raw http response for streaming
// https://developer.twitter.com/en/docs/tweets/filter-realtime/api-reference/post-statuses-filter
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
// https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline
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
	err = json.Unmarshal(resBytes, &output)
	if err != nil {
		return []StatusesUserTimelineOutput{}, err
	}

	return output, nil
}

func (c Client) AccountsScheduledTweetsPost(input AccountsScheduledTweetsInput, accountId string) (*http.Response, error) {
	uri := "https://ads-api.twitter.com/8/accounts/" + accountId + "/scheduled_tweets"
	params := processParams(input)
	res, err := c.executeRequest(http.MethodPost, uri, params)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	log.Println("you've posted a scheduled tweet.", res.Body)
	return res, nil
}
