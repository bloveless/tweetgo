# Tweetgo

This library is not meant to be an exhaustive API for Twitter, yet. My goal is two-fold. First, I want it to be
dead simple to add new endpoints to this library. Second I want the benefit of strong types that we all enjoy
in Go. I've modeled this library after aws-sdk-go where the input to every endpoint is a well-defined struct,
and the output from every endpoint is also a well-defined struct.

A quick example of posting a tweet. You'll have to provide the consumer key, consumer secret, access token, and access
token secret yourself. The endpoints for getting the oauth request token and converting that into and oauth access token
have been implemented as OAuthRequestTokenGet and OAuthAccessTokenGet. In case you want to have your application
automatically authorize the user and retrieve access tokens on their behalf. Checkout the provided examples.

The easiest way to run the examples is to cd into the `example` directory and run `go run .`

```go
package main

func main() {
    tc := tweetgo.Client{
        OAuthConsumerKey:       "OAuthConsumerKey",
        OAuthConsumerSecret:    "OAuthConsumerSecret",
        OAuthAccessToken:       "OAuthAccessToken",
        OAuthAccessTokenSecret: "OAuthAccessTokenSecret",
    }

    input := tweetgo.StatusesUpdateInput{
        Status: "Hello Ladies + Gentlemen, a signed OAuth request!",
    }

    output, err := tc.StatusesUpdatePost(input)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v", output)
}
```

If you look in the model.go you'll see that every option for updating a status is in the struct StatusesUpdateInput
and all the possible fields have been documented in the StatusesUpdateOutput.

If you follow these steps you can add any endpoints that you need easily and give back to the community!

## Setup for local development

If you are using Go mod in your project you can add something like the following:

```
replace github.com/bloveless/tweetgo => ../tweetgo
```

to your go.mod file before the require lines (replace ../tweetgo with the path to your local copy of tweetgo).

I learned about that command from this blog post. [Using "replace" in go.mod to point to your local module](https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/)

## How to add a new endpoint

The second goal for this project is to make it super easy to add new endpoints. These are the steps you'll need to
follow in order to add another endpoint. I'll be adding the GET [statuses/user_timeline](https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline)
endpoint.

### Step 1: Create the input model

Refer to the link above to generate the input model. Below is the input I created for the statuses/user_timeline
endpoint. Twitter users form encoded input when receiving a request and responds in JSON. This doesn't change much
about implementing a new endpoint other than you'll use the "schema" tag when you are writing the input structs, and the
"json" tag when you are writing the output structs. Every input field needs to be a pointer so we can handle null
values. There are helper function for converting between standard types and pointers to make this a little easier.
Checkout `value.go` to see what I'm talking about. If you add a type that hasn't been used before, it will be helpful
for you to add a conversion func to `value.go`.

```go
type StatusesUserTimelineInput struct {
	UserID         *int64  `schema:"user_id"`
	ScreenName     *string `schema:"screen_name"`
	SinceID        *int64  `schema:"since_id"`
	Count          *int    `schema:"count"`
	MaxID          *int64  `schema:"max_id"`
	TrimUser       *bool   `schema:"trim_user"`
	ExcludeReplies *bool   `schema:"exclude_replies"`
	IncludeRts     *bool   `schema:"include_rts"`
}
```

### Step 2: Create the output model

Refer to the link above to generate the output model. There are only a few different types of returns from Twitter API's
so there are some utility structs that you can compose responses from. They are unexported and are intended to be
composed within other structs. The statuses user timeline endpoint returns a list of tweets so the output struct is very
simple.

```go
type StatusesUserTimelineOutput struct {
	tweet
}
```

### Step 3: Create endpoint code

All of the request signing and execution has been wrapped up into utility functions which you can find in `request.go`.
Hopefully this will make adding new endpoints extremely simple. For this endpoint I added the following code to 
`client.go`.

```go
func (c Client) StatusesUserTimelineGet(input StatusesUserTimelineInput) ([]StatusesUserTimelineOutput, error) { // 1) change to the correct input/output types here
	uri := "https://api.twitter.com/1.1/statuses/user_timeline.json" // 2) Change to the correct URI here
	params := processParams(input)

	res, err := c.executeRequest(http.MethodGet, uri, params) // 3) Change to the correct http method here
	if err != nil {
		return []StatusesUserTimelineOutput{}, err // 4) Change to the correct output type here
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body) // NOTE: This can be changed to parse form value output
	if err != nil {
		return []StatusesUserTimelineOutput{}, err // 5) Change to the correct output type here
	}

	var output []StatusesUserTimelineOutput // 6) Change to the correct output type here
	json.Unmarshal(resBytes, &output)

	return output, nil
}
```

Follow the six steps above, and you'll have your own endpoint up and running.

NOTE: As far as I can tell twitter uses JSON output on nearly all of it's endpoints. The oauth endpoints use form value
rather than JSON output so support functions exist for processing that type of data as well. Look at OAuthRequestTokenGet
which uses gorilla/scheme for parsing the data into the output format.

There you go! You've added your own endpoint!
