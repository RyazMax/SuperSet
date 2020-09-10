package static

import "net/http"

const (
	// StaticRootDir dir where css and js located
	StaticRootDir = "src/static"
)

// Handler serves static for web
func Handler() http.Handler {
	return http.FileServer(http.Dir(StaticRootDir))
}
