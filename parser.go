package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strings"
	"time"
)

type Feed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title string `xml:"title"`
	Link  struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Thumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
}

// Performs the parsing
func getFeedEntries(u string) ([]Entry, error) {
	u = u + ".rss"
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	// simulate request from a browser to avoid risk of being blocked.
	request.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Mobile Safari/537.36",
	)

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	feedByte, _ := ioutil.ReadAll(response.Body)
	var feed Feed
	err = xml.Unmarshal(feedByte, &feed)
	if err != nil {
		return nil, err
	}

	return feed.Entries, nil
}

func isValidSubredditURL(rawURL string) bool {
	url, err := netUrl.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	validHostname := "reddit.com"
	if url.Host != validHostname && url.Host != "www."+validHostname {
		return false
	}

	unCleanPath := strings.Split(url.Path, "/")
	if len(unCleanPath) > 4 {
		return false
	}

	var cleanPath []string
	for _, str := range unCleanPath {
		if str != "" {
			cleanPath = append(cleanPath, str)
		}
	}
	pathLen := len(cleanPath)
	return pathLen > 0 && cleanPath[0] == "r" && pathLen <= 2
}
