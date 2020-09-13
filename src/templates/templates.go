package templates

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"../models"
	"../universe"
)

const (
	// HTMLRootDir points to folder where html templates are located
	HTMLRootDir = "src/templates/html"
)

// Templates
var (
	IndexPageTemplate   *template.Template
	ProfilePageTemplate *template.Template
	LoginPageTemplate   *template.Template
	NewUserPageTemplate *template.Template
	LabelPageTemplate   *template.Template
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

	LoginPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "login.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}

	NewUserPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "new_user.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}

	LabelPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "label.html")}, baseTemplates...)...)
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

// LoginHandler accepts "/login" requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("Login")
		password := r.FormValue("Password")
		sess, err := universe.Get().Auth.Login(login, password)
		if err != nil {
			log.Println(sess)
			w.Write([]byte("Can't login"))
			return
		}
		if sess == nil {
			w.Write([]byte("Can't login"))
			return
		}
		addAuthCookie(w, sess.ID)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		err := LoginPageTemplate.Execute(w, struct{}{})
		if err != nil {
			log.Println(err)
		}
	}
}

// NewUserHandler accepts "/new_user" requests
func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("Login")
		password := r.FormValue("Password")
		// Add validate
		sess, err := universe.Get().Auth.NewUser(&models.User{Login: login, PasswordShadow: password})
		if err != nil {
			log.Println(err)
			w.Write([]byte("Can't create User"))
			return
		}
		addAuthCookie(w, sess.ID)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		err := NewUserPageTemplate.Execute(w, struct{}{})
		if err != nil {
			log.Println(err)
		}
	}
}

// LogoutHandler accepts "/logout" requests
// Login required
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		log.Println(err)
		return
	}
	err = universe.Get().Auth.Logout(cookie.Value)
	if err != nil {
		log.Println(err)
	}
	unsetAuthCookie(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// LabelHandler accepts "/label" requests leading to labeling page
// Login required
func LabelHandler(w http.ResponseWriter, r *http.Request) {
	err := LabelPageTemplate.Execute(w, struct{}{})
	if err != nil {
		log.Println(err)
	}
}

// Handler returns http.Handler that serves web-site routes
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/new_user", notLoginRequired(http.HandlerFunc(NewUserHandler)))
	mux.Handle("/login", notLoginRequired(http.HandlerFunc(LoginHandler)))

	mux.Handle("/profile", loginRequired(http.HandlerFunc(ProfilePageHandler)))
	mux.Handle("/label", loginRequired(http.HandlerFunc(LabelHandler)))
	mux.Handle("/logout", loginRequired(http.HandlerFunc(LogoutHandler)))
	mux.Handle("/", passUserName(http.HandlerFunc(IndexPageHandler)))

	return mux
}
