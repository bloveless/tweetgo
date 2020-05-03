package main

import (
	"encoding/json"
	"fmt"
	"io"
	"tweetgo"
)

func postStatus(c config, status string) {
	tc := getTwitterClient(c)

	input := tweetgo.StatusesUpdateInput{
		Status: tweetgo.String(status),
	}

	output, err := tc.StatusesUpdatePost(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%+v\n", output)
}

func streamTweets(c config, hashtag string) {
	fmt.Println("Beginning to stream kubernetes tweets")
	tc := getTwitterClient(c)

	input := tweetgo.StatusesFilterInput{
		Track: tweetgo.String(hashtag),
	}

	output, err := tc.StatusesFilterPostRaw(input)
	if err != nil {
		panic(err)
	}

	for {
		tweet := tweetgo.StatusesFilterOutput{}
		err := json.NewDecoder(output.Body).Decode(&tweet)
		if err == io.EOF {
			fmt.Println("End of file")
		}

		if err != nil {
			panic(err)
		}

		fmt.Printf("%#v\n\n\n", tweet)
	}
}
