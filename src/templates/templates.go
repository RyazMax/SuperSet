package templates

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"../models"
	"../project_manager"
	"../universe"
)

const (
	// HTMLRootDir points to folder where html templates are located
	HTMLRootDir = "src/templates/html"
	// MediaDir dir where uploaded datasets are located
	MediaDir = "data/media"
)

// Templates
var (
	IndexPageTemplate      *template.Template
	ProfilePageTemplate    *template.Template
	LoginPageTemplate      *template.Template
	NewUserPageTemplate    *template.Template
	LabelPageTemplate      *template.Template
	NewProjectPageTemplate *template.Template
	UploadPageTemplate     *template.Template
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

	NewProjectPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "new_project.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}

	UploadPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "upload.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}
}

// IndexPageHandler handles "/"" requests
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	err := IndexPageTemplate.Execute(w, createDataOnContext(r.Context()))
	if err != nil {
		log.Println(err)
	}
}

// ProfilePageHandler handles "/profile" requests
func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	err := ProfilePageTemplate.Execute(w, createDataOnContext(r.Context()))
	if err != nil {
		log.Println(err)
	}
}

// LoginHandler accepts "/login" requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		login := r.FormValue("Login")
		password := r.FormValue("Password")
		sess, err := universe.Get().Auth.Login(login, password)
		if err != nil {
			log.Println(sess)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if sess == nil {
			data["InputError"] = "Введен неверный логин/пароль"
			err := LoginPageTemplate.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
			return
		}
		addAuthCookie(w, sess.ID)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := LoginPageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}

}

// NewUserHandler accepts "/new_user" requests
func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		login := r.FormValue("Login")
		password := r.FormValue("Password")
		// Add validate
		sess, err := universe.Get().Auth.NewUser(&models.User{Login: login, PasswordShadow: password})
		if err != nil {
			log.Println(err)
			data["InputError"] = "Пользователь с таким именем уже существует"
			err := NewUserPageTemplate.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
			return
		}
		addAuthCookie(w, sess.ID)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := NewUserPageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
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
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		return
	}

	task, err := universe.Get().TaskManager.TakeTask([]string{"Tester"})
	if err != nil {
		log.Println("LabelHandler error on getting task", err)
	}
	if task != nil {
		data["TaskFound"] = true
		data["InputType"] = task.Schema.InputSchema.InputType()
		// some proccesing for task Data
		data["OutputType"] = task.Schema.OutputSchema.OutputType()
		// some processing for task schema
	}
	err = LabelPageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// NewProjectHandler accepts "/new_project" requests leading to new project construction
// Login required
func NewProjectHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		prName := r.FormValue("projectName")
		ischemaType := r.FormValue("ischema")
		oschemaType := r.FormValue("oschema")
		ischema := models.InputTypeToStructMap(ischemaType)
		oschema := models.OutputTypeToStructMap(oschemaType)
		if ischema == nil || oschema == nil {
			data["InputError"] = "Неверный формат ввода/вывода"
			return
		}

		uname, _ := data["UserName"].(string)
		user, _ := universe.Get().UserRepo.GetByLogin(uname)
		prAggr := project_manager.ProjectAggr{
			Project: models.Project{
				Name:    prName,
				OwnerID: int(user.ID),
			},
			Schema: models.ProjectSchema{
				InputSchema:  ischema,
				OutputSchema: oschema,
			},
		}
		ok, err := universe.Get().ProjectManager.Create(&prAggr)
		if !ok {
			log.Println(err)
		} else if err != nil {
			data["InputError"] = err.Error()
		}
		http.Redirect(w, r, "/project", http.StatusTemporaryRedirect)
		return
	}
	err := NewProjectPageTemplate.Execute(w, createDataOnContext(r.Context()))
	if err != nil {
		log.Println(err)
	}
}

// UploadHandler handles upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		r.ParseMultipartForm(10 << 20)
		fmt.Println(r.RequestURI)
		m := r.MultipartForm
		projName := r.FormValue("projectName")

		proj, err := universe.Get().ProjectRepo.GetByName(projName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if proj == nil {
			// handle nil project
		}
		schema, err := universe.Get().SchemaRepo.GetByProjectID(proj.ID)
		if err != nil || schema == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//get the *fileheaders
		files := m.File["dataFile"]
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			task, err := schema.InputSchema.Validate(files[i].Filename, file)
			if err != nil {
				// handler validate error
				log.Println(err)
				return
			}
			task.ProjectID = proj.ID

			id, err := universe.Get().TaskManager.PutTask(projName, task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Do something with id
			savename := schema.InputSchema.SaveName(files[i].Filename, id)
			if savename != "" {
				//create destination file making sure the path is writeable.
				log.Println("Saving", files[i].Filename)
				dst, err := os.Create(path.Join(MediaDir, savename))
				defer dst.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//copy the uploaded file to the destination file
				if _, err := io.Copy(dst, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		fmt.Fprintf(w, "Successfully uploaded")
	} else {
		err := UploadPageTemplate.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

// Handler returns http.Handler that serves web-site routes
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/new_user", notLoginRequired(http.HandlerFunc(NewUserHandler)))
	mux.Handle("/login", notLoginRequired(http.HandlerFunc(LoginHandler)))

	mux.Handle("/new_project", loginRequired(http.HandlerFunc(NewProjectHandler)))
	mux.Handle("/profile", loginRequired(http.HandlerFunc(ProfilePageHandler)))
	mux.Handle("/label", loginRequired(http.HandlerFunc(LabelHandler)))
	mux.Handle("/logout", loginRequired(http.HandlerFunc(LogoutHandler)))
	mux.Handle("/upload", loginRequired(http.HandlerFunc(UploadHandler)))
	mux.Handle("/", passUserName(http.HandlerFunc(IndexPageHandler)))

	return mux
}
