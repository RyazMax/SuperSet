package templates

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><head><title>Set</title></head><body><h1>SuperSet</h1></body></html>"))
}

// Handler returns http.Handler that serves web-site routes
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	return mux
}
