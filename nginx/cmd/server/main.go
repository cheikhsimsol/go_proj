package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc(
		"/api/v1",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Request sent to %s\n", r.URL.Path)

			fmt.Fprintln(w, "Welcome to the protected home page!")
		},
	)

	fmt.Println("Listenning on port 5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}

}
