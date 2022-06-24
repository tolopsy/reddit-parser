package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", ":9000", "Server port to listen to")
	flag.Parse()

	http.HandleFunc("/", getParsedSubredditData)

	fmt.Println("Service running in port", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
