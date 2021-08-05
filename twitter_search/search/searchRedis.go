package search

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"twitter_search/docs"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type SearchController struct {
	rClient redis.Conn
	eClient *elasticsearch.Client
}

func NewSearchController(r *gin.Engine, eClient *elasticsearch.Client, rClient redis.Conn) {
	searchController := SearchController{
		rClient: rClient,
		eClient: eClient,
	}
	enterpriseRouter := r.Group("/v1")
	docs.SwaggerInfo.Title = "Title"
	enterpriseRouter.GET("/tweet/search", searchController.GetTweets)
	enterpriseRouter.GET("/tweet/searchByLang", searchController.GetTweetsByLang)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}

func (ctrl *SearchController) GetDocId(search string) []MapRow {
	sep := " "
	keys := strings.Split(strings.ToLower(search), sep)
	docMap := make(map[string]int)
	for _, key := range keys {
		value, err := redis.Strings(ctrl.rClient.Do("LRANGE", key, 0, -1))
		if err != nil {
			log.Fatal(err)
		}
		for _, docId := range value {
			if _, ok := docMap[docId]; !ok {
				docMap[docId] = 1
			} else {
				docMap[docId] += 1
			}
		}

	}

	return SortMap(docMap)
}

// @Summary Get raw tweet based on search
// @Description Search tweets
// @ID get-tweets
// @Accept json
// @Produce json
// @Param search query string true "search"
// @Success 200 {object} SearchTweetResponse
// @Tags Tweets
// @Router /v1/tweet/search [get]
func (ctrl *SearchController) GetTweets(c *gin.Context) {
	// default past 1 month
	search := c.Query("search")
	finalOutput := &SearchTweetResponse{}

	docIds := ctrl.GetDocId(search)
	for _, id := range docIds {
		tweetObj := &EsTweetObj{
			TweetId: id.Key,
			Tweet:   &TweetStruct{},
		}
		request := esapi.GetRequest{Index: "tweets", DocumentID: id.Key}
		response, err := request.Do(context.Background(), ctrl.eClient)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorRs{
				Code:    http.StatusBadRequest,
				Message: "Cannot find the keyword",
			})
			return
		}
		var results map[string]interface{}
		json.NewDecoder(response.Body).Decode(&results)
		//println(results)
		// json.Unmarshal(results, &tweet)
		if results["_source"] != nil {
			esTweet := results["_source"].(map[string]interface{})
			//mapstructure.WeakDecode(esTweet, &tweetObj)
			// println(esTweet["Tweet"].(map[string]interface{})["text"].(string))
			tweetObj.Tweet.Name = esTweet["Tweet"].(map[string]interface{})["user"].(map[string]interface{})["name"].(string)
			tweetObj.Tweet.CreatedAt = esTweet["Tweet"].(map[string]interface{})["created_at"].(string)
			tweetObj.Tweet.Lang = esTweet["Tweet"].(map[string]interface{})["lang"].(string)
			tweetObj.Tweet.Text = esTweet["Tweet"].(map[string]interface{})["text"].(string)
			finalOutput.Tweets = append(finalOutput.Tweets, *tweetObj)
		}
	}
	c.JSON(http.StatusOK, finalOutput)
}

// @Summary Get raw tweet based on search
// @Description Search tweets by lang, accepted language:
// @Description English (default)	en
// @Description  Arabic	ar
// @Description  Bengali	bn
// @Description  Czech	cs
// @Description  Danish	da
// @Description  German	de
// @Description  Greek	el
// @Description  Spanish	es
// @Description  Persian	fa
// @Description  Finnish	fi
// @Description  Filipino	fil
// @Description  French	fr
// @Description  Hebrew	he
// @Description  Hindi	hi
// @Description  Hungarian	hu
// @Description  Indonesian	id
// @Description  Italian	it
// @Description  Japanese	ja
// @Description  Korean	ko
// @Description  Malay	msa
// @Description  Dutch	nl
// @Description  Norwegian	no
// @Description  Polish	pl
// @Description  Portuguese	pt
// @Description  Romanian	ro
// @Description  Russian	ru
// @Description  Swedish	sv
// @Description  Thai	th
// @Description  Turkish	tr
// @Description  Ukrainian	uk
// @Description  Urdu	ur
// @Description  Vietnamese	vi
// @Description  Chinese (Simplified)	zh-cn
// @Description  Chinese (Traditional)	zh-tw
// @ID get-tweets-byLang
// @Accept json
// @Produce json
// @Param search query string true "search by lang"
// @Success 200 {object} SearchTweetResponse
// @Tags Tweets
// @Router /v1/tweet/searchByLang [get]
func (ctrl *SearchController) GetTweetsByLang(c *gin.Context) {
	// default past 1 month
	language := c.Query("search")
	finalOutput := &SearchTweetResponse{}
	var (
		buffer bytes.Buffer
	)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"Tweet.lang": language,
			},
		},
	}
	if err := json.NewEncoder(&buffer).Encode(query); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorRs{
			Code:    http.StatusInternalServerError,
			Message: "Error encoding the query",
		})
		return
	}
	response, err := ctrl.eClient.Search(ctrl.eClient.Search.WithIndex("tweets"), ctrl.eClient.Search.WithBody(&buffer))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, ErrorRs{
			Code:    http.StatusBadGateway,
			Message: "Error getting data from elasticSearch",
		})
		return
	}
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		tweetObj := &EsTweetObj{
			TweetId: "",
			Tweet:   &TweetStruct{},
		}
		esTweet := hit.(map[string]interface{})["_source"].(map[string]interface{})
		tweetObj.TweetId = esTweet["TweetId"].(string)
		tweetObj.Tweet.Name = esTweet["Tweet"].(map[string]interface{})["user"].(map[string]interface{})["name"].(string)
		tweetObj.Tweet.CreatedAt = esTweet["Tweet"].(map[string]interface{})["created_at"].(string)
		tweetObj.Tweet.Lang = esTweet["Tweet"].(map[string]interface{})["lang"].(string)
		tweetObj.Tweet.Text = esTweet["Tweet"].(map[string]interface{})["text"].(string)
		finalOutput.Tweets = append(finalOutput.Tweets, *tweetObj)
	}
	c.JSON(http.StatusOK, finalOutput)
}
