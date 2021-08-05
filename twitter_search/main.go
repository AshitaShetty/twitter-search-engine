package main

import (
	"log"
	"twitter_search/search"

	"github.com/elastic/go-elasticsearch"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	eClient, redisClient := CreateClients()
	search.NewSearchController(router, eClient, redisClient)
}

func CreateClients() (*elasticsearch.Client, redis.Conn) {
	// set elasticSearch
	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// set redis
	rclient, err := redis.Dial("tcp", ":6379")
	if err != nil {
		rclient.Close()
	}

	return esClient, rclient
}
