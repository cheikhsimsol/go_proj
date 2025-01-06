package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/hashicorp/consul/api"
)

func registerHandler(
	serviceId string,
	service *api.AgentService,
) {

	//address := fmt.Sprintf("")
	pathname := fmt.Sprintf("/%s", serviceId)
	address := fmt.Sprintf("%s:%v", service.Address, service.Port)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {

			url := url.URL{
				Scheme: "http",
				Host:   address,
				Path:   strings.TrimPrefix(r.In.URL.Path, pathname),
			}

			log.Println("forwarding:", url.Path, "-->", url.String())

			r.SetXForwarded()
			r.Out.URL = &url

		},
	}

	http.HandleFunc(
		pathname+"/",
		proxy.ServeHTTP,
	)
}

func main() {

	consulAddress := "http://localhost:8500"

	config := &api.Config{
		Address: consulAddress, // Consul address
	}

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	filter := `Tags contains "gateway"`
	services, err := client.Agent().ServicesWithFilter(filter)
	if err != nil {
		log.Fatalf("Failed to fetch filtered services: %v", err)
	}

	// Print the filtered services
	for serviceId, service := range services {

		registerHandler(
			serviceId,
			service,
		)

		fmt.Printf("Proxying path: /%s  Address: %s, Port: %d, Tags: %v\n",
			serviceId, service.Address, service.Port, service.Tags)

	}

	http.ListenAndServe(":7000", nil)
}
