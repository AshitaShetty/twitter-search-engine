package twitter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type TweetData struct {
	Tweet        string
	LikeCount    int
	RetweetCount int
}

func GetCredentials() Credentials { // can set the variable using export, but have hardcoded it for simplicity
	creds := Credentials{
		ConsumerKey:       "****************",
		ConsumerSecret:    "*****************",
		AccessToken:       "*******************",
		AccessTokenSecret: "********************",
	}

	return creds
}

func GetClient(creds *Credentials) (*twitter.Client, error) {

	// These values are the API key and API key secret
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// These values are the consumer access token and consumer access token secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verify := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err := client.Accounts.VerifyCredentials(verify)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// print out the Twitter handle of the account we have used to authenticate with
	fmt.Println("Successfully authenticated using the following account : ", user.ScreenName)
	return client, nil
}

// func SearchTweets(client *twitter.Client) error {

// 	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
// 		Query: "Trump",
// 		Lang:  "en",
// 	})

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	//fmt.Printf("length %d", len(search.Statuses))
// 	for _, v := range search.Statuses {
// 		tweet := TweetData{
// 			Tweet:        v.Text,
// 			LikeCount:    v.FavoriteCount,
// 			RetweetCount: v.RetweetCount,
// 		}
// 		fmt.Printf("%+v\n", tweet)
// 	}
// 	return nil
// }

func StreamTweets(client *twitter.Client) error {
	params := &twitter.StreamFilterParams{
		Track:         []string{"covid"},
		StallWarnings: twitter.Bool(true),
		//Language:      []string{"en",""},
	}
	stream, err := client.Streams.Filter(params)
	if err != nil {
		fmt.Println(err)
		return err
	}
	demux := twitter.NewSwitchDemux()

	//set kafka_producer
	broker := "localhost"
	topic := "myKafkaTopic"

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)
	fmt.Printf("Starting to stream..")
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
		marshalTweet, _ := json.Marshal(tweet)
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          marshalTweet,
			Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		}, nil)
	}
	for message := range stream.Messages {
		demux.Handle(message)
	}
	return nil
}
