package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	redditBaseUrl = "https://www.reddit.com"
)

type AutoCompleteResp struct {
	Children []struct {
		Data struct {
			DisplayName string `json:"display_name"`
		} `json:"data"`
	} `json:"children"`
}

type PostsResponse struct {
	Children []struct {
		Data struct {
			Url      string `json:"url"`
			IsVideo  bool   `json:"is_video"`
			Over18   bool   `json:"over_18"`
			PostHint string `json:"post_hint"`
		} `json:"data"`
	} `json:"children"`
}

func RedditAutoComplete(query string, over18 bool) (*AutoCompleteResp, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/subreddit_autocomplete_v2.json?query=%s&include_over_18=%t", redditBaseUrl, query, over18),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch autocomplete: %s", err)
	}

	req.Header.Set("user-agent", "AlwaysAdd/0.0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch autocomplete: %s", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode > 299 {
		return nil, errors.New("Non 200 status code")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var autoCompleteResp struct {
		Data AutoCompleteResp `json:"data"`
	}
	json.Unmarshal(bodyBytes, &autoCompleteResp)

	return &autoCompleteResp.Data, nil
}

func RedditGetPost(subreddit string, over18 bool) (*PostsResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/r/%s.json?include_over_18=%t", redditBaseUrl, subreddit, over18),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %s", err)
	}

	req.Header.Set("user-agent", "AlwaysAdd/0.0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %s", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode > 299 {
		return nil, errors.New("Non 200 status code")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var postsResponse struct {
		Data PostsResponse `json:"data"`
	}
	json.Unmarshal(bodyBytes, &postsResponse)

	return &postsResponse.Data, nil
}

func RedditSaveImage(url string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %s", err)
	}

	req.Header.Set("user-agent", "AlwaysAdd/0.0.1")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %s", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode > 299 {
		return "", errors.New("Non 200 status code")
	}

	splitUrl := strings.Split(url, ".")
	ext := splitUrl[len(splitUrl)-1]

	fileName := fmt.Sprintf("reddit_image.%s", ext)
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("Failed to create image %s: %s", fileName, err)
	}

	defer file.Close()
	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return "", fmt.Errorf("Failed to copy image %s: %s", fileName, err)
	}

	return fileName, nil
}
