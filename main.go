package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	customTransport = http.DefaultTransport
)

func andleRequest(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		log.Printf("Error creating proxy request: %v", err)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Set(name, value)
		}
	}

	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		log.Printf("Error forwarding request: %v", err)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

// func main() {
// 	remote, err := url.Parse("https://www.youtube.com/")
// 	if err != nil {
// 		panic(err)
// 	}
// 	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			log.Println(r.URL)
// 			r.Host = remote.Host
// 			w.Header().Set("X-Ben", "Rad")
// 			p.ServeHTTP(w, r)
// 		}
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(remote)
// 	proxy.Transport = &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	http.HandleFunc("/", handler(proxy))
// 	err = http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// }

var (
	// URL de base du serveur de destination
	baseURL *url.URL
)

func init() {
	// Spécifiez l'URL de base ici
	var err error
	baseURL, err = url.Parse("https://www.youtube.com") // Remplacez par votre URL de base
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Construire l'URL de la requête de destination
	targetURL := baseURL.ResolveReference(r.URL)

	// Vérification que l'URL de destination est bien formée
	if targetURL.Scheme == "" || targetURL.Host == "" {
		http.Error(w, "Invalid proxy target URL", http.StatusInternalServerError)
		log.Printf("Invalid proxy target URL: %s", targetURL.String())
		return
	}

	log.Printf("Forwarding request to: %s", targetURL.String())

	// Créer une nouvelle requête HTTP pour le serveur de destination
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		log.Printf("Error creating proxy request: %v", err)
		return
	}

	// Copier les en-têtes de la requête originale
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Set(name, value)
		}
	}

	// Envoyer la requête au serveur de destination
	resp, err := http.DefaultTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		log.Printf("Error forwarding request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copier les en-têtes de la réponse
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Set(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	// Copier le corps de la réponse (flux vidéo)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

func main() {
	http.HandleFunc("/", handleRequest)

	log.Println("Starting proxy server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting proxy server: %v", err)
	}
}
