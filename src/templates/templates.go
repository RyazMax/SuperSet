package templates

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

const (
	// HTMLRootDir points to folder where html templates are located
	HTMLRootDir = "src/templates/html"
)

// Templates
var (
	IndexPageTemplate   *template.Template
	ProfilePageTemplate *template.Template
)

// Loading templates
func init() {
	baseTemplates := []string{
		//	path.Join(HTMLRootDir, "index.html"),
		//	path.Join(HTMLRootDir, "profile.html"),
		path.Join(HTMLRootDir, "base.html"),
		path.Join(HTMLRootDir, "header.html"),
		path.Join(HTMLRootDir, "footer.html"),
	}

	var err error
	IndexPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "index.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}

	ProfilePageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "profile.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}
}

// IndexPageHandler handles "/"" requests
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	err := IndexPageTemplate.Execute(w, struct{ UserName string }{name})
	if err != nil {
		log.Println(err)
	}
}

// ProfilePageHandler handles "/profile" requests
func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	err := ProfilePageTemplate.Execute(w, struct{}{})
	if err != nil {
		log.Println(err)
	}
}

// Handler returns http.Handler that serves web-site routes
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/profile", ProfilePageHandler)
	mux.HandleFunc("/", IndexPageHandler)

	return mux
}
