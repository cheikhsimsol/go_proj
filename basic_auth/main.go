package main

import (
	"fmt"
	"net/http"
)

// handler for the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username != "admi" || password != "password" {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Secret API\"")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Unauthorized")
		return
	}

	fmt.Fprintln(w, "Welcome to the protected home page!")
}

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
