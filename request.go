package twitter

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

func processParams(input interface{}) url.Values {
	encoder := schema.NewEncoder()

	params := url.Values{}
	encoder.Encode(input, params)

	for k := range params {
		if len(params.Get(k)) == 0 {
			params.Del(k)
		}
	}

	return params
}

func (c Client) executeRequest(method, uri string, params url.Values) (*http.Response, error) {
	req, err := c.getSignedRequest(method, uri, params)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		errBody, _ := ioutil.ReadAll(res.Body)
		return nil, errors.New(string(errBody))
	}

	return res, nil
}

func bodyToValues(body io.ReadCloser) (url.Values, error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return url.Values{}, err
	}

	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		return url.Values{}, err
	}

	return values, nil
}

func (c Client) getSignedRequest(method, uri string, params url.Values) (*http.Request, error) {
	nonce := generateNonce()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	sr := signatureRequest{
		method:    method,
		uri:       uri,
		nonce:     nonce,
		timestamp: timestamp,
		params:    params,
	}

	oauthSignature := c.signature(sr)

	hp := headerParameters{
		oauthNonce:     nonce,
		oauthSignature: oauthSignature,
		oauthTimestamp: timestamp,
	}

	authHeader := c.getOauthAuthorizationHeader(hp)

	req, err := http.NewRequest(sr.method, sr.uri, strings.NewReader(sr.params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authHeader)

	return req, nil
}

func generateNonce() string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 48)
	for i := range b {
		b[i] = allowed[rand.Intn(len(allowed))]
	}
	return string(b)
}

type signatureRequest struct {
	method    string
	uri       string
	nonce     string
	timestamp string
	params    url.Values
}

func (c Client) signature(sr signatureRequest) string {
	values := url.Values{}
	values.Add("oauth_consumer_key", c.OAuthConsumerKey)
	values.Add("oauth_nonce", sr.nonce)
	values.Add("oauth_signature_method", "HMAC-SHA1")
	values.Add("oauth_timestamp", sr.timestamp)
	values.Add("oauth_version", "1.0")

	if c.OAuthAccessToken != "" {
		values.Add("oauth_token", c.OAuthAccessToken)
	}

	for k := range sr.params {
		values.Add(url.QueryEscape(k), sr.params.Get(k))
	}

	parameterString := strings.ReplaceAll(values.Encode(), "+", "%20")

	signatureBaseString := strings.ToUpper(sr.method) +
		"&" + url.QueryEscape(strings.Split(sr.uri, "?")[0]) +
		"&" + url.QueryEscape(parameterString)

	signingKey := url.QueryEscape(c.OAuthConsumerSecret) + "&" + url.QueryEscape(c.OAuthAccessTokenSecret)

	return calculateSignature(signatureBaseString, signingKey)
}

func calculateSignature(base, key string) string {
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(base))
	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

type headerParameters struct {
	oauthNonce     string
	oauthSignature string
	oauthTimestamp string
}

func (c Client) getOauthAuthorizationHeader(p headerParameters) string {
	authHeader := "OAuth " +
		"oauth_consumer_key=\"" + url.QueryEscape(c.OAuthConsumerKey) + "\", " +
		"oauth_nonce=\"" + url.QueryEscape(p.oauthNonce) + "\", " +
		"oauth_signature=\"" + url.QueryEscape(p.oauthSignature) + "\", " +
		"oauth_signature_method=\"HMAC-SHA1\", " +
		"oauth_timestamp=\"" + url.QueryEscape(p.oauthTimestamp) + "\", "

	if c.OAuthAccessToken != "" {
		authHeader += "oauth_token=\"" + url.QueryEscape(c.OAuthAccessToken) + "\", "
	}

	authHeader += "oauth_version=\"1.0\""

	return authHeader
}
