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
		fmt.Errorf("could not load configuration: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	for _, resource := range config.Resources {
		url, _ := url.Parse(resource.Destination_URL)
		proxy := NewProxy(url)
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(proxy, url, resource.Endpoint))
	}
	if err := http.ListenAndServe(config.Server.Host+":"+config.Server.Listen_port, mux); err != nil {
		return fmt.Errorf("could not start the server: %v", err)
	}
	return nil
}
