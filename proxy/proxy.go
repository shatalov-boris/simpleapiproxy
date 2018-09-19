package apiproxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Set the proxied request's host to the destination host (instead of the
// source host).  e.g. http://foo.com proxying to http://bar.com will ensure
// that the proxied requests appear to be coming from http://bar.com
//
// For both this function and queryCombiner (below), we'll be wrapping a
// Handler with our own HandlerFunc so that we can do some intermediate work
func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

// Allow cross origin resource sharing
func addHeaders(handler http.Handler, corsDomain string, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		r.Header.Set("x-api-key", apiKey)
		handler.ServeHTTP(w, r)
	})
}

// Combine the two functions above with http.NewSingleHostReverseProxy
func Proxy(remoteURL string, corsDomain string, apiKey string) http.Handler {
	// pull the root url we're proxying to from an environment variable.
	serverURL, err := url.Parse(remoteURL)

	if err != nil {
		log.Fatal("URL failed to parse")
	}

	// initialize our reverse proxy
	reverseProxy := httputil.NewSingleHostReverseProxy(serverURL)
	// wrap that proxy with our sameHost function
	singleHosted := sameHost(reverseProxy)
	// and finally allow CORS
	return addHeaders(singleHosted, corsDomain, apiKey)
}
