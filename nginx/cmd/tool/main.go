package main

import (
	"log"
	"os"
	"text/template"
)

type Service struct {
	Path string
	Host string
}

type Config struct {
	Services []Service // A slice of services to define multiple proxy locations
}

func main() {
	// Path to the template file
	templatePath := "nginx.conf.tmpl"

	// Define the configuration to be rendered
	config := []Service{
		{Path: "/api/", Host: "http://localhost:5000"},
		{Path: "/auth/", Host: "http://localhost:4000"},
	}

	// Load the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Error loading template: %v", err)
	}

	// Create the output file
	outputFile, err := os.Create("nginx.conf")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Render the template with the configuration
	err = tmpl.Execute(outputFile, config)
	if err != nil {
		log.Fatalf("Error rendering template: %v", err)
	}

	log.Println("nginx.conf has been generated successfully.")
}
