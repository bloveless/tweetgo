package tests

import (
	"fmt"
	"github.com/bloveless/tweetgo"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockClient struct {
	requests []*http.Request
	responses []*http.Response
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)

	if len(m.responses) > 0 {
		response := m.responses[0]
		m.responses = m.responses[1:]

		return response, nil
	}

	return nil, nil
}

func TestCanProcessParamsAndOmitNilValues(t *testing.T) {
	c := tweetgo.Client{
		HTTPClient:
	}
	type testStruct struct {
		TestString     *string  `schema:"test_string"`
		TestBool       *bool    `schema:"test_bool"`
		TestInt        *int     `schema:"test_int"`
		TestInt64      *int64   `schema:"test_int64"`
		TestFloat64    *float64 `schema:"test_float64"`
		TestNilString  *string  `schema:"test_nil_string"`
		TestNilBool    *bool    `schema:"test_nil_bool"`
		TestNilInt     *int     `schema:"test_nil_int"`
		TestNilInt64   *int64   `schema:"test_nil_int64"`
		TestNilFloat64 *float64 `schema:"test_nil_float64"`
	}

	ts := testStruct{
		TestString:     tweetgo.String("test"),
		TestBool:       tweetgo.Bool(true),
		TestInt:        tweetgo.Int(10),
		TestInt64:      tweetgo.Int64(20),
		TestFloat64:    tweetgo.Float64(3.49),
		TestNilString:  nil,
		TestNilBool:    nil,
		TestNilInt:     nil,
		TestNilInt64:   nil,
		TestNilFloat64: nil,
	}

	o := tweetgo.processParams(ts)

	expected := url.Values{
		"test_string":  {"test"},
		"test_bool":    {"true"},
		"test_int":     {"10"},
		"test_int64":   {"20"},
		"test_float64": {"3.49"},
	}

	if !reflect.DeepEqual(expected, o) {
		fmt.Printf("Expected: %+v\n", expected)
		fmt.Printf("Actual:   %+v\n", o)
		t.Fail()
	}
}

func TestCanProcessStructsThatContainZeroValues(t *testing.T) {
	type testZeroStruct struct {
		TestZeroString  *string  `schema:"test_zero_string"`
		TestZeroBool    *bool    `schema:"test_zero_bool"`
		TestZeroInt     *int     `schema:"test_zero_int"`
		TestZeroInt64   *int64   `schema:"test_zero_int64"`
		TestZeroFloat64 *float64 `schema:"test_zero_float64"`
	}

	ts := testZeroStruct{
		TestZeroString:  tweetgo.String(""),
		TestZeroBool:    tweetgo.Bool(false),
		TestZeroInt:     tweetgo.Int(0),
		TestZeroInt64:   tweetgo.Int64(0),
		TestZeroFloat64: tweetgo.Float64(0.00),
	}

	o := tweetgo.processParams(ts)

	expected := url.Values{
		"test_zero_string":  {""},
		"test_zero_bool":    {"false"},
		"test_zero_int":     {"0"},
		"test_zero_int64":   {"0"},
		"test_zero_float64": {"0"},
	}

	if !reflect.DeepEqual(expected, o) {
		fmt.Printf("Expected: %+v\n", expected)
		fmt.Printf("Actual:   %+v\n", o)
		t.Fail()
	}
}

// Using one of the twitter examples make sure we can correctly calculate the signature
// https://developer.twitter.com/en/docs/basics/authentication/oauth-1-0a/creating-a-signature
func TestCorrectlyCalculatesSignatureForStatusesUpdate(t *testing.T) {
	tc := tweetgo.Client{
		OAuthConsumerKey:       "xvz1evFS4wEEPTGEFPHBog",
		OAuthConsumerSecret:    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OAuthAccessToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		OAuthAccessTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	sr := tweetgo.signatureRequest{
		method:    http.MethodPost,
		uri:       "https://api.twitter.com/1.1/statuses/update.json",
		nonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		timestamp: "1318622958",
		params: url.Values{
			"status":           {"Hello Ladies + Gentlemen, a signed OAuth request!"},
			"include_entities": {"true"},
		},
	}

	sig, err := tc.signature(sr)
	if err != nil {
		t.Fatalf("Signature generation failed: %s", err.Error())
	}

	expected := "hCtSmYh+iHYCEqBWrE7C7hYmtUk="

	if sig != expected {
		t.Fatalf("sig: %s != expected: %s", sig, expected)
	}
}

func TestCorrectlyCalculatesSignatureForStatusesUpdateWithGetParameters(t *testing.T) {
	tc := tweetgo.Client{
		OAuthConsumerKey:       "xvz1evFS4wEEPTGEFPHBog",
		OAuthConsumerSecret:    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OAuthAccessToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		OAuthAccessTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	sr := tweetgo.signatureRequest{
		method:    http.MethodPost,
		uri:       "https://api.twitter.com/1.1/statuses/update.json?include_entities=true",
		nonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		timestamp: "1318622958",
		params: url.Values{
			"status": {"Hello Ladies + Gentlemen, a signed OAuth request!"},
		},
	}

	sig, err := tc.signature(sr)
	if err != nil {
		t.Fatalf("Signature generation failed: %s", err.Error())
	}

	expected := "hCtSmYh+iHYCEqBWrE7C7hYmtUk="

	if sig != expected {
		t.Fatalf("sig: %s != expected: %s", sig, expected)
	}
}

// Duplicated params will change the signature
func TestCorrectlyCalculatesSignatureForStatusesUpdateWithDuplicatedGetParameters(t *testing.T) {
	tc := tweetgo.Client{
		OAuthConsumerKey:       "xvz1evFS4wEEPTGEFPHBog",
		OAuthConsumerSecret:    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		OAuthAccessToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		OAuthAccessTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	sr := tweetgo.signatureRequest{
		method:    http.MethodPost,
		uri:       "https://api.twitter.com/1.1/statuses/update.json?include_entities=true",
		nonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		timestamp: "1318622958",
		params: url.Values{
			"status":           {"Hello Ladies + Gentlemen, a signed OAuth request!"},
			"include_entities": {"true"},
		},
	}

	sig, err := tc.signature(sr)
	if err != nil {
		t.Fatalf("Signature generation failed: %s", err.Error())
	}

	expected := "p8ht/l/ns5JbCNn8mP+TsRgp4U0="

	if sig != expected {
		t.Fatalf("sig: %s != expected: %s", sig, expected)
	}
}
