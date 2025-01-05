package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Service configuration
	serviceID := "example-service-2"
	serviceName := "example-service-2"
	servicePort := 8090
	serviceAddress := getLocalIP()
	consulAddress := "http://localhost:8500"

	// Start a simple HTTP server
	go startHTTPServer(servicePort)

	// Register the service with Consul
	err := registerServiceWithConsul(serviceID, serviceName, serviceAddress, consulAddress, servicePort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Printf("Service %s is running on %s:%d and registered with Consul.\n", serviceName, serviceAddress, servicePort)

	// Prevent the program from exiting
	select {}
}

// registerServiceWithConsul registers the service with Consul
func registerServiceWithConsul(
	serviceID,
	serviceName,
	serviceAddress,
	consulAddress string,
	servicePort int,
) error {
	// Create a Consul client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("failed to create Consul client: %v", err)
	}

	// Define the service registration
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", serviceAddress, servicePort),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	// Register the service
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	return nil
}

// startHTTPServer starts a simple HTTP server with a health check endpoint
func startHTTPServer(port int) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Printf("Starting HTTP server on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// getLocalIP gets the local IP address of the machine
func getLocalIP() string {
	// Fallback to localhost if environment variables are not set
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	return host
}
