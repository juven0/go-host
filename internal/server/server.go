package server

import (
	"fmt"
	"go-host/internal/configs"
	"net/http"
	"net/url"
)

func Run() error {
	config, err := configs.NewConfiguration()
	if err != nil {
		return fmt.Errorf("could not load configuration: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	for _, resource := range config.Resources {
		url, err := url.Parse(resource.Destination_URL)
		if err != nil {
			return fmt.Errorf("invalid destination URL for endpoint %s: %v", resource.Endpoint, err)
		}

		proxy := NewProxy(url)
		if proxy == nil {
			return fmt.Errorf("failed to create proxy for endpoint %s", resource.Endpoint)
		}

		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(proxy, url, resource.Endpoint))
	}

	serverAddr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Listen_port)
	fmt.Printf("Starting server on %s\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}

	return nil
}
