package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strings"
	"time"
)

type feed struct {
	Entries []entry `xml:"entry" json:"entry"`
}

type entry struct {
	Title string `xml:"title" json:"title"`

	Author struct {
		Name string `xml:"name" json:"name"`
		URI  string `xml:"uri" json:"uri"`
	} `xml:"author" json:"author"`

	Content struct {
		Type  string `xml:"type,attr" json:"type"`
		Value string `xml:",innerxml" json:"content"`
	} `xml:"content" json:"content"`

	Link struct {
		Href string `xml:"href,attr" json:"href"`
	} `xml:"link" json:"link"`

	Thumbnail struct {
		URL string `xml:"url,attr" json:"url"`
	} `xml:"thumbnail" json:"thumbnail"`

	Updated   time.Time `xml:"updated" json:"updated"`
	Published time.Time `xml:"published" json:"published"`
}

// Performs the parsing
func getFeedEntries(u string) ([]entry, error) {
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
	var feed feed
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

	// confirm host
	if url.Host != validHostname && url.Host != "www."+validHostname {
		return false
	}

	if !strings.HasSuffix(url.Path, "/") {
		url.Path = url.Path + "/"
	}
	path := strings.Split(url.Path, "/")
	// confirm url path is within the formats ["/r", "/r/", "/r/topic", "/r/topic/"]
	// this means the path can not be /r/topic/another-string
	return len(path) > 1 && path[1] == "r" && len(path) <= 4
}
