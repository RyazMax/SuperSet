package templates

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"../models"
	"../project_manager"
	"../universe"
)

const (
	// HTMLRootDir points to folder where html templates are located
	HTMLRootDir = "src/templates/html"
	// MediaDir dir where uploaded datasets are located
	MediaDir = "src/static/media"
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
	StartLabelPageTemplate *template.Template
	DownloadPageTemplate   *template.Template
	ProjectPageTemplate    *template.Template
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

	StartLabelPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "start_label.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}

	DownloadPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "download.html")}, baseTemplates...)...)
	if err != nil {
		log.Fatal(err)
	}

	ProjectPageTemplate, err = template.ParseFiles(append([]string{path.Join(HTMLRootDir, "project.html")}, baseTemplates...)...)
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
	data := createDataOnContext(r.Context())
	uname, _ := data["UserName"].(string)
	user, _ := universe.Get().UserRepo.GetByLogin(uname)
	if user == nil {
		log.Printf("Not found user with name %s", user.Login)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Type == models.AdminUser {
		projects, err := universe.Get().ProjectRepo.GetByOwnerID(int(user.ID))
		if err != nil {
			log.Println("[ProfilePageHandler]", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data["Projects"] = projects
	} else if user.Type == models.RegularUser {
		http.Error(w, "Not implemented yet", http.StatusNotImplemented)
		return
	} else {
		log.Printf("User %s has invalid type %d", user.Login, user.Type)
		http.Error(w, "Incorrect user type", http.StatusBadRequest)
		return
	}

	err := ProfilePageTemplate.Execute(w, data)
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
		log.Println("BEGIN")
		pid, _ := strconv.Atoi(r.FormValue("projectID"))
		tid, _ := strconv.Atoi(r.FormValue("queuetaskID"))
		oid, _ := strconv.Atoi(r.FormValue("originID"))
		log.Println("BEFORE GET PROJECT", pid, tid)
		schema, err := universe.Get().SchemaRepo.GetByProjectID(pid)
		log.Println("SCHEMA", schema)
		if err != nil || schema == nil {
			http.Error(w, "ERROR or nil schema", http.StatusInternalServerError)
			return
		}
		ltask, err := schema.OutputSchema.FormatLabeledTask(r)
		if err != nil {
			log.Println(err)
			// handle err
			return
		}
		ltask.OriginID = oid

		log.Println("Smth", ltask)
		tsk := models.TaskAggr{ID: tid, Tsk: models.Task{ProjectID: pid}}
		err = universe.Get().TaskManager.LabelTask(&tsk, ltask)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	projects := strings.Split(r.FormValue("Projects"), ";")
	log.Println(projects)
	task, err := universe.Get().TaskManager.TakeTask(projects)
	if err != nil {
		log.Println("LabelHandler error on getting task", err)
	}
	data["Projects"] = strings.Join(projects, ";")
	if task != nil {
		data["TaskFound"] = true
		data["ProjectID"] = task.Tsk.Tsk.ProjectID
		data["OriginID"] = task.Tsk.Tsk.ID
		data["TaskID"] = task.Tsk.ID // For queue
		data["InputType"] = task.Schema.InputSchema.InputType()
		data["OutputType"] = task.Schema.OutputSchema.OutputType()
		task.Schema.InputSchema.FormatInputData(&task.Tsk.Tsk, data)
		task.Schema.OutputSchema.FormatOutputData(&task.Tsk.Tsk, data)
	}

	// For tpr
	/*data["InputType"] = "TableData"
	data["OutputType"] = "PlainText"
	data["TaskFound"] = true
	data["ColNames"] = []string{"Длина чашелистника", "Ширина чашелистника", "Длина лепестка"}
	data["ColVals"] = []string{"5.1", "3.5", "1.4"}*/
	err = LabelPageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// StartLabelHandler starts labeling process
func StartLabelHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		projName := r.FormValue("projectName")
		ok, err := universe.Get().ProjectManager.CheckGrant(projName, data["UserName"].(string))
		if ok {
			http.Redirect(w, r, "/label?Projects="+strings.Join([]string{projName}, ";"), 303) // Redirect
		} else {
			data["Error"] = err.Error()
		}
	}
	err := StartLabelPageTemplate.Execute(w, data)
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
		http.Redirect(w, r, fmt.Sprintf("/project?name=%s", prName), http.StatusTemporaryRedirect)
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

		ok, err := universe.Get().ProjectManager.IsOwner(data["UserName"].(string), projName)
		if !ok {
			data["Error"] = err.Error()
			err := UploadPageTemplate.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
			return
		}
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

		mapping := []interface{}{}
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
			mapping = append(mapping, []interface{}{files[i].Filename, id})
		}
		w.Header().Set("Content-Disposition", "attachment; filename=mapping.json")
		encoder := json.NewEncoder(w)
		encoder.Encode(mapping)
		return
	} else {
		err := UploadPageTemplate.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	if r.Method == "POST" {
		projectName := r.FormValue("projectName")
		outputType := r.FormValue("type")
		originID := r.FormValue("originID")

		ok, err := universe.Get().ProjectManager.IsOwner(data["UserName"].(string), projectName)
		if !ok {
			data["Error"] = err.Error()
			err := DownloadPageTemplate.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
			return
		}
		project, _ := universe.Get().ProjectRepo.GetByName(projectName)

		var ltsks []models.LabeledTask
		if originID != "" {
			origin, err := strconv.Atoi(originID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			tsk, err := universe.Get().LabeledRepo.GetByOriginID(project.ID, origin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if tsk != nil {
				ltsks = []models.LabeledTask{*tsk}
			}
		} else {
			ltsks, err = universe.Get().LabeledRepo.GetByProjectID(project.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Disposition", "attachment; filename=dataset."+outputType)
		if outputType == "json" {
			jsonDataFormatter(w, ltsks)
		} else if outputType == "csv" {
			csvDataFormatter(w, ltsks)
		} else {
			http.Error(w, "Bad format", http.StatusNotImplemented)
		}
		return
	}
	err := DownloadPageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// ProjectPageHandler handles requests for project page
func ProjectPageHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())

	addError := r.FormValue("adderror")
	if addError != "" {
		data["AddError"] = addError
	}
	deleteError := r.FormValue("deleteerror")
	if deleteError != "" {
		data["DeleteError"] = deleteError
	}

	projectName := r.FormValue("name")
	project, err := universe.Get().ProjectRepo.GetByName(projectName)
	if err != nil {
		log.Printf("[ProjectPageHandler] Error on getting project %s, %v", projectName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Project"] = project
	// Add .formatInfo for schemas
	schema, err := universe.Get().SchemaRepo.GetByProjectID(project.ID)
	if err != nil {
		log.Printf("[ProjectPageHandler] Error on getting schema %d, %v", project.ID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Schema"] = schema

	err = ProjectPageTemplate.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

// GrantHandler do requests on "/grant"
func GrantHandler(w http.ResponseWriter, r *http.Request) {
	data := createDataOnContext(r.Context())
	isDel := r.FormValue("delete") != ""
	var addError, deleteError string
	if !isDel {
		ok, err := universe.Get().ProjectManager.AddGrant(data["UserName"].(string),
			r.FormValue("name"),
			r.FormValue("username"))
		if !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err != nil {
			addError = err.Error()
		}
	} else if isDel {
		ok, err := universe.Get().ProjectManager.DeleteGrant(data["UserName"].(string),
			r.FormValue("name"),
			r.FormValue("username"))
		if !ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err != nil {
			deleteError = err.Error()
		}
	} else {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	projName := r.FormValue("name")
	http.Redirect(w, r, fmt.Sprintf("/project?name=%s&adderror=%s&deleteerror=%s", projName, addError, deleteError), 303)
	return
}

// Handler returns http.Handler that serves web-site routes
func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/new_user", notLoginRequired(http.HandlerFunc(NewUserHandler)))
	mux.Handle("/login", notLoginRequired(http.HandlerFunc(LoginHandler)))

	mux.Handle("/new_project", loginRequired(http.HandlerFunc(NewProjectHandler)))
	mux.Handle("/profile", loginRequired(http.HandlerFunc(ProfilePageHandler)))
	mux.Handle("/project", loginRequired(http.HandlerFunc(ProjectPageHandler)))
	mux.Handle("/start_label", loginRequired(http.HandlerFunc(StartLabelHandler)))
	mux.Handle("/label", loginRequired(http.HandlerFunc(LabelHandler)))
	mux.Handle("/logout", loginRequired(http.HandlerFunc(LogoutHandler)))
	mux.Handle("/upload", loginRequired(http.HandlerFunc(UploadHandler)))
	mux.Handle("/download", loginRequired(http.HandlerFunc(DownloadHandler)))
	mux.Handle("/grant", loginRequired(http.HandlerFunc(GrantHandler)))
	mux.Handle("/", passUserName(http.HandlerFunc(IndexPageHandler)))

	//return recoverMiddleware(mux)
	return mux
}
