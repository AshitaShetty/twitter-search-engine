# twitter-search-engine
Search through streaming tweets using a keyword/language

# Description:
Tweets are streamed in real time into a kafka topic which is then consumed by the consumer and then fed to elasticSearch. The search API is hosted on swagger which can search tweets by keywords or by language.

# Getting Started

## Installation on Windows

Download and extract ZooKeeper from http://zookeeper.apache.org/releases.html

Download and extract Kafka from http://kafka.apache.org/downloads.html

Download and extract Elasticsearch from https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.2.4.zip

Older versions of redis are available on windows (which is enough for this project) follow the instruictions on https://github.com/microsoftarchive/redis/releases


## Start Kafka service

Open cmd, navigate to folder containing Zookeeper and start the zookeeper using the command

+ `$ .\bin\zkserver`

Open another cmd, navigate to folder containing Kafka and start the kafka service using the command
+ `$ .\bin\windows\kafka-server-start.bat .\config\server.properties`

## Start ElasticSearch

Open cmd, navigate to folder containing elasticsearch and use command

+ `$ .\bin\elasticsearch.bat`

## Start Redis

Nagivate to the folder containing redis (64bit) and run redis-server.exe

## Get data from twitter

To get data from twitter, get authentication code from  https://developer.twitter.com/en/apply-for-access and update program twitter_producer.go with the credentials.

# How to run the service

Open three instances of go to run the producer,consumer and search seperately.

## twitter_producer
This folder contains all the files required to start the twitter streaming service and add the tweets into Kafka topic.

To run this producer, run the following command :

+ `$ go run main.go`

## twitter_consumer
This folder contains all the files which are required to start the consumtion of the tweets and add them to elasticSearch.

Each tweet is parsed as a json and added to the elasticsearch. An inverse index is created based on the tweet.Text by removing stop words and special characters. This inverse index is stored as a key-list pair in redis,which will later be used in the search service

To run this consumer, run the following command :

+ `$ go run main.go`

## twitter-search
This folder contains all the files required to start the twitter Search service and host the API on Swagger.

User is able to search the tweets by keyword and language

To run this search engine, run the following command :

+ `$ go run main.go`

The API is hosted on : http://localhost:8080/swagger/index.html#/