package api

import "net/http"

// Handler return http.Handler to serve service requests
func Handler() http.Handler {
	apiMux := http.NewServeMux()
	return apiMux
}
