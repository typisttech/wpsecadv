// TODO: Use npx serve instead.

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is empty, defaulting to 10080")
		port = "10080"
	}
	log.Printf("Port:\t%s", port) //gosec:disable G706

	fixture := os.Getenv("FEED_FIXTURE")
	if fixture == "" {
		log.Fatal("FEED_FIXTURE is empty")
	}
	log.Printf("Fixture:\t%s", fixture) //gosec:disable G706

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fixture)
	})

	log.Printf("Serving fixture on http://localhost:%s", port) //gosec:disable G706

	if err := http.ListenAndServe(":"+port, nil); err != nil { //gosec:disable G114
		log.Fatal(err)
	}
}
