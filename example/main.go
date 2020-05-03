package main

// In order to run these examples create an example/config.json file with the following contents:
//
// {
//    "oauth_consumer_key":"XXXXXXXXXXXXXXXXXXXXXXXXX",
//    "oauth_consumer_secret":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
// }
//
// Just put in your real twitter consumer key and secret. The auth function will fill in the access key and secret
// then the update status function will post a status update to your twitter account
func main() {
	c := auth()

	postStatus(c, "Hello Ladies + Gentlemen, a signed OAuth request!")

	streamTweets(c, "twitter")
}
