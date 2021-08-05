package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"twitter_consumer/indexing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/garyburd/redigo/redis"
)

type EsTweetObj struct {
	TweetId string      `json:"TweetId"`
	Tweet   interface{} `json:"Tweet"`
}

func SetConsumer() {

	//setting up consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"myKafkaTopic"}, nil)

	// setting up elastic search
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	//setting up redis for index storage
	client, err := redis.Dial("tcp", ":6379")
	if err != nil {
		client.Close()
	}

	filepath := "./indexing/stop_words.txt"
	tweetId := 1
	tweetIndex := make(map[string][]int)

	for {
		msg, err := c.ReadMessage(-1)
		var tweet twitter.Tweet
		if err == nil {
			json.Unmarshal(msg.Value, &tweet)
			//fmt.Printf("Message on %s: %s\n", msg.TopicPartition, tweet.Text)
			stopWords := indexing.LoadStopWordsFromFile(filepath)
			finalTweet := indexing.RemoveStopWords(stopWords, tweet.Text)
			CreateIndex(tweetId, tweetIndex, finalTweet, client)
			LoadElasticSearch(tweetId, tweet, es)
			tweetId++
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func CreateIndex(tweetId int, tweetIndex map[string][]int, tweets []string, client redis.Conn) {
	for _, word := range tweets {
		client.Do("RPUSH", word, tweetId)
	}
}

func LoadElasticSearch(tweetId int, tweet interface{}, es *elasticsearch.Client) {

	fmt.Printf("Message on elasticS %s\n", tweet.(twitter.Tweet).Text)

	id := strconv.Itoa(tweetId)
	fmt.Println("Document id:" + id)
	esTweet := EsTweetObj{id, tweet}
	jsonTweet, _ := json.Marshal(esTweet)
	request := esapi.IndexRequest{Index: "tweets", DocumentID: id, Body: strings.NewReader(string(jsonTweet))}
	_, err := request.Do(context.Background(), es)
	if err != nil {
		println(err)
	}
}
