package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/atanurdemir/gatekeeper/src/service"
	"github.com/atanurdemir/gatekeeper/src/types"
	"github.com/atanurdemir/gatekeeper/src/utils"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Read and buffer the body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Extract the user ID from the request context
	userID, err := utils.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized: UserID not found in context", http.StatusUnauthorized)
		return
	}

	// Prepare RequestContext without the Body
	ctx := types.RequestContext{
		Path:   r.URL.Path,
		Method: r.Method,
		IP:     utils.GetIP(r),
		UserID: userID,
	}

	// Process the request
	result, err := service.Check.ProcessRequest(r.Context(), ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !result.Status {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(result)
		return
	}

	// Reassign the buffered body to the request for forwarding
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Set up and execute the reverse proxy
	setupAndExecuteProxy(w, r, bodyBytes)
}

func setupAndExecuteProxy(w http.ResponseWriter, r *http.Request, bodyBytes []byte) {
	destinationURL, err := extractDestinationURL(r)
	if err != nil {
		http.Error(w, "Invalid destination URL", http.StatusBadRequest)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(destinationURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = destinationURL.Scheme
		req.URL.Host = destinationURL.Host
		req.URL.Path = destinationURL.Path
		req.URL.RawQuery = destinationURL.RawQuery
		req.Host = destinationURL.Host
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("proxying request to %s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path+"?"+req.URL.RawQuery)
	}

	proxy.ServeHTTP(w, r)
}

func extractDestinationURL(r *http.Request) (*url.URL, error) {
	destination := r.Header.Get("X-Destination-URL")
	return url.Parse(destination)
}
