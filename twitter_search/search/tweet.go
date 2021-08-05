package search

import "sort"

type SearchTweetResponse struct {
	Tweets []EsTweetObj `json:"tweets"`
}

type EsTweetObj struct {
	TweetId string       `json:"TweetId"`
	Tweet   *TweetStruct `json:"Tweet"`
}

type TweetStruct struct {
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	Text      string `json:"text"`
}

type ErrorRs struct {

	// code
	Code int64 `json:"code,omitempty"`

	// message
	// Required: true
	Message string `json:"message"`
}

type MapRow struct {
	Key   string
	Value int
}

func SortMap(m map[string]int) []MapRow {
	var sortedDocMap []MapRow
	for k, v := range m {
		sortedDocMap = append(sortedDocMap, MapRow{k, v})
	}

	sort.Slice(sortedDocMap, func(i, j int) bool {
		return sortedDocMap[i].Value > sortedDocMap[j].Value
	})
	return sortedDocMap
}
