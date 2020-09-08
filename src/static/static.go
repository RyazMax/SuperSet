package static

import "net/http"

// Handler serves static for web
func Handler() http.Handler {
	return http.FileServer(http.Dir("/tmp"))
}
