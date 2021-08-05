package indexing

import (
	"io/ioutil"
	"strings"
	"twitter_consumer/util"
)

func LoadStopWordsFromFile(filePath string) map[string]string {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	//fmt.Print(string(b))
	return LoadStopWordsFromString(string(b))
}

func LoadStopWordsFromString(wordlist string) map[string]string {
	sep := "\r\n"
	words := strings.Split(strings.ToLower(wordlist), sep)
	// fmt.Print(words)
	stopWordsMap := make(map[string]string)

	for _, word := range words {
		stopWordsMap[word] = ""
	}
	return stopWordsMap
}

func RemoveStopWords(stopWordsMap map[string]string, tweetText string) []string {
	var result []string
	sep := " "
	tweet := util.RemoveSpecialChar(tweetText)
	words := strings.Split(strings.ToLower(tweet), sep)
	for _, word := range words {
		if _, ok := stopWordsMap[string(word)]; ok {
			result = append(result, "")
		} else {
			result = append(result, word)
			//result = append(result, " ")
		}
	}
	return result
}
