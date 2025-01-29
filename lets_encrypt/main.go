package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, secure world with Let's Encrypt!")
}

func main() {
	// Replace this with your actual domain name
	domain := os.Getenv("DOMAIN_NAME")

	// Create an autocert manager
	certManager := autocert.Manager{
		Cache:      autocert.DirCache(".certs"), // Folder for storing certificates
		Prompt:     autocert.AcceptTOS,          // Automatically agree to the TOS
		HostPolicy: autocert.HostWhitelist(domain),
	}

	// Configure the HTTPS server
	server := &http.Server{
		Addr:      ":443", // HTTPS port
		Handler:   http.HandlerFunc(helloHandler),
		TLSConfig: certManager.TLSConfig(), // Use autocert TLS config
	}

	// Start a background HTTP server for Let's Encrypt HTTP-01 challenge
	go func() {
		log.Println("Starting HTTP server for Let's Encrypt challenge...")
		err := http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		if err != nil {
			log.Fatalf("HTTP server failed: %s\n", err)
		}
	}()

	// Start the HTTPS server
	log.Printf("Starting HTTPS server on https://%s...\n", domain)
	err := server.ListenAndServeTLS("", "") // TLS config is handled by autocert
	if err != nil {
		log.Fatalf("HTTPS server failed: %s\n", err)
	}
}
