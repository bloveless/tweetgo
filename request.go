package tweetgo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type requestMaker interface {
	Do(req *http.Request) (*http.Response, error)
}

type nonceMaker interface {
	Generate() string
}

type currentTimer interface {
	GetCurrentTime() int64
}

func processParams(input interface{}) url.Values {
	v := reflect.ValueOf(input)

	params := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Tag.Get("schema")
		field := v.Field(i)

		// Convert to non-pointer version
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		// skip unset/invalid fields
		if !field.IsValid() {
			continue
		}

		// get the actual value
		value := field.Interface()

		if value != nil {
			// convert to string based on underlying type
			switch value.(type) {
			case string:
				params.Add(name, value.(string))
			case bool:
				params.Add(name, strconv.FormatBool(value.(bool)))
			case int:
				params.Add(name, strconv.FormatInt(int64(value.(int)), 10))
			case int64:
				params.Add(name, strconv.FormatInt(value.(int64), 10))
			case float64:
				params.Add(name, strconv.FormatFloat(value.(float64), 'f', -1, 64))
			}
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

	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		return nil, errors.New("Status: " + res.Status + " - Body: " + string(b))
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
	nonce := c.Noncer.Generate()
	timestamp := strconv.FormatInt(c.Timer.GetCurrentTime(), 10)

	sr := signatureRequest{
		method:    method,
		uri:       uri,
		nonce:     nonce,
		timestamp: timestamp,
		params:    params,
	}

	oauthSignature, err := c.signature(sr)
	if err != nil {
		return nil, err
	}

	hp := headerParameters{
		oauthNonce:     nonce,
		oauthSignature: oauthSignature,
		oauthTimestamp: timestamp,
	}

	authHeader := c.getOauthAuthorizationHeader(hp)

	req, err := http.NewRequest(sr.method, sr.uri, strings.NewReader(sr.params.Encode()))

	if method == http.MethodGet {
		u, err := url.Parse(sr.uri)
		if err != nil {
			return nil, err
		}

		u.RawQuery = sr.params.Encode()

		req, err = http.NewRequest(sr.method, u.String(), nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authHeader)

	return req, nil
}

type signatureRequest struct {
	method    string
	uri       string
	nonce     string
	timestamp string
	params    url.Values
}

func (c Client) signature(sr signatureRequest) (string, error) {
	uri, err := url.Parse(sr.uri)
	if err != nil {
		return "", err
	}

	values := uri.Query()

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

	return calculateSignature(signatureBaseString, signingKey), nil
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
