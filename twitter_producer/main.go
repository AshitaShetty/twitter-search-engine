package main

import (
	"fmt"
	"os"
	"twitter_producer/twitter"
)

func main() {
	creds := twitter.GetCredentials()
	client, err := twitter.GetClient(&creds)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	twitter.StreamTweets(client)
}
