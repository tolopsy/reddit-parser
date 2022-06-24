package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getParsedSubredditData(w http.ResponseWriter, r *http.Request) {
	var url string
	switch r.Method {
	case http.MethodGet:
		url = r.URL.Query().Get("url")
	case http.MethodPost:
		var payload struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request payload: %s", err.Error()), http.StatusBadRequest)
			return
		}
		url = payload.URL
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isValidSubredditURL(url) {
		http.Error(w, fmt.Sprintln("Invalid subreddit url:", url), http.StatusBadRequest)
		return
	}
	entries, err := getFeedEntries(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&entries)
}
