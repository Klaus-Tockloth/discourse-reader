package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var progressBar = "\rRetrieving page %d ..."

/*
getQueryData gets general data for site.
*/
func getQueryData() ([]byte, error) {
	requestURL := *query
	status, bodyData, err := getDiscourseData(requestURL)
	if err != nil {
		return nil, fmt.Errorf("error [%v] requesting general query data", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status [%v] requesting general query data (%v)", status, string(bodyData))
	}

	return bodyData, nil
}

// DiscourseCategory represents minimal response data for further processing.
type DiscourseCategory struct {
	TopicList struct {
		MoreTopicsURL string `json:"more_topics_url"`
	} `json:"topic_list"`
}

/*
getCategoryData gets data (list of topics) for category.
*/
func getCategoryData() (string, error) {
	var tmpPages []string
	var categoryPages string
	pageNumber := 1

	requestURL := fmt.Sprintf("https://%s/c/communities/-/%d/l/latest.json?ascending=false", *forum, *category)
	for i := 0; i < *pages; i++ {
		discourseCategory := DiscourseCategory{}
		fmt.Printf(progressBar, pageNumber)
		pageNumber++
		status, bodyData, err := getDiscourseData(requestURL)
		if err != nil {
			return "", fmt.Errorf("error [%v] requesting category data", err)
		}

		if status != http.StatusOK {
			return string(bodyData), fmt.Errorf("unexpected HTTP status [%v] requesting category data (%v)", status, string(bodyData))
		}

		err = json.Unmarshal(bodyData, &discourseCategory)
		if err != nil {
			return "", fmt.Errorf("error [%v] unmarshaling category data", err)
		}

		tmpPages = append(tmpPages, string(bodyData))

		if discourseCategory.TopicList.MoreTopicsURL == "" {
			break
		}

		requestURL = fmt.Sprintf("https://%s%s", *forum, discourseCategory.TopicList.MoreTopicsURL)
		time.Sleep(time.Duration(*sleeptime) * time.Second) // we won't performed this action too many times (see user rate limits)
	}

	fmt.Printf("\n")

	// build JSON category pages array
	categoryPages = fmt.Sprintf("{ \"category_pages\": [\n%s\n] }\n", strings.Join(tmpPages, ",\n"))

	return categoryPages, nil
}

// DiscourseTopic represents minimal response data (of first chunk) for further processing.
type DiscourseTopic struct {
	PostStream struct {
		Stream []int `json:"stream"`
	} `json:"post_stream"`
}

/*
getTopicData gets data (list of posts) for topic.

first read  : meta data for topic (e.g. 'stream' with post ids) + first chunk (e.g. 20) of posts
second read : 50 posts (posts from first read are redundant included)
third read  : 50 posts
...

E.g. topic has 3 posts:
first read  : meta data for topic + 3 posts
second read : 3 posts (posts from first read are redundant included)

E.g. topic has 143 posts:
first read  : meta data for topic + 20 of posts
second read : 50 posts (posts from first read are redundant included)
third read  : 50 posts
fourth read : 43 posts
*/
func getTopicData() (string, error) {
	var discourseTopic DiscourseTopic
	var tmpPages []string
	var topicPages string
	pageNumber := 1

	// first request contains meta data of topic (with first chunks of posts and complete list of post ids)
	fmt.Printf(progressBar, pageNumber)
	pageNumber++
	requestURL := fmt.Sprintf("https://%s/t/-/%d.json", *forum, *topic)
	status, bodyData, err := getDiscourseData(requestURL)
	if err != nil {
		return "", fmt.Errorf("error [%v] requesting topic data", err)
	}

	if status != http.StatusOK {
		return string(bodyData), fmt.Errorf("unexpected HTTP status [%v] requesting topic data (%v)", status, string(bodyData))
	}

	err = json.Unmarshal(bodyData, &discourseTopic)
	if err != nil {
		return "", fmt.Errorf("error [%v] unmarshaling first chunk of topic data", err)
	}

	metaData := string(bodyData)

	// process list of post ids (in pages of 50 posts)
	chunkSize := 50
	requestURL = fmt.Sprintf("https://%s/t/%d/posts.json?", *forum, *topic)
	i := 0
	streamID := 0
	for i, streamID = range discourseTopic.PostStream.Stream {
		// add post ids from (posts) stream to request URL
		if i == 0 {
			requestURL += fmt.Sprintf("post_ids[]=%d", streamID)
		} else {
			requestURL += fmt.Sprintf("&post_ids[]=%d", streamID)
		}

		if (i+1)%chunkSize == 0 {
			fmt.Printf(progressBar, pageNumber)
			pageNumber++
			time.Sleep(time.Duration(*sleeptime) * time.Second) // we won't performed this action too many times (see user rate limits)
			status, bodyData, err = getDiscourseData(requestURL)
			if err != nil {
				return "", fmt.Errorf("error [%v] requesting topic data", err)
			}

			if status != http.StatusOK {
				return "", fmt.Errorf("unexpected HTTP status [%v] requesting topic data (%v)", status, string(bodyData))
			}

			tmpPages = append(tmpPages, string(bodyData))
			requestURL = fmt.Sprintf("https://%s/t/%d/posts.json?", *forum, *topic)
		}
	}

	// read last chunk (only if something is left)
	if (i+1)%chunkSize != 0 {
		fmt.Printf(progressBar, pageNumber)
		time.Sleep(time.Duration(*sleeptime) * time.Second) // we won't performed this action too many times (see user rate limits)
		status, bodyData, err = getDiscourseData(requestURL)
		if err != nil {
			return "", fmt.Errorf("error [%v] requesting topic data", err)
		}

		if status != http.StatusOK {
			return "", fmt.Errorf("unexpected HTTP status [%v] requesting topic data (%v)", status, bodyData)
		}

		tmpPages = append(tmpPages, string(bodyData))
	}

	fmt.Printf("\n")

	// build JSON category pages array
	topicPages = fmt.Sprintf("{\n"+
		"\"meta_data\": %s,\n"+
		"\"post_data\": [\n%s\n]\n"+
		"}\n", metaData, strings.Join(tmpPages, ",\n"))

	return topicPages, nil
}
