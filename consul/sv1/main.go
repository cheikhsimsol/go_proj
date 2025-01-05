package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Service configuration
	serviceID := "example-service-1"
	serviceName := "example-service"
	servicePort := 8080
	serviceAddress := "127.0.0.1"
	consulAddress := "http://localhost:8500"

	config := &api.Config{
		Address: consulAddress, // Consul address
	}

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	// Start a simple HTTP server
	go startHTTPServer(servicePort)

	// Register the service with Consul
	err = registerServiceWithConsul(
		serviceID, serviceName, serviceAddress, servicePort, client)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Printf("Service %s is running on %s:%d and registered with Consul.\n", serviceName, serviceAddress, servicePort)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigs

	// Deregister the service upon shutdown
	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Fatalf("Failed to deregister service: %v", err)
	}
	fmt.Println("Service deregistered from Consul")

}

// registerServiceWithConsul registers the service with Consul
func registerServiceWithConsul(
	serviceID,
	serviceName,
	serviceAddress string,
	servicePort int,
	consulClient *api.Client,
) error {
	// Create a Consul client

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
	err := consulClient.Agent().ServiceRegister(registration)
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
