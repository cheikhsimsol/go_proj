package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const spaDir = "./dist" // Change this to your SPA directory

func spaHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(spaDir, r.URL.Path)

	// If the request has no file extension, serve index.html
	if !strings.Contains(filepath.Base(path), ".") {
		http.ServeFile(w, r, filepath.Join(spaDir, "index.html"))
		return
	}

	// Serve static files
	http.FileServer(http.Dir(spaDir)).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", spaHandler)

	log.Println("Serving SPA on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
