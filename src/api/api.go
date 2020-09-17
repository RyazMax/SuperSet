package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
 Dataset loading api:

 POST /api/data - upload data
 JSON:
	project: id of project
	token: user token
	type: some type
	count: unsigned
	data: some data

 GET /api/data - download labeled data
 Queryparams:
	 project: id of project
	 type: json|csv
	 token: user token
	 [limit]: unsigned| 0 - for limitless
	 [origin]: id of origin task
	 [ts]: timestamp to find greater
*/

func GetDataHandler(w http.ResponseWriter, r *http.Request) {

}

func PostDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	fmt.Println(r.RequestURI)
	file, handler, err := r.FormFile("dataFile")
	if err != nil {
		fmt.Println("Error with uploading file", err)
		return
	}
	defer file.Close()

	fmt.Println(handler.Filename, handler.Header)

	tempFile, err := ioutil.TempFile("images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	fmt.Fprintf(w, "Successfully uploaded")
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PostDataHandler(w, r)
	} else if r.Method == http.MethodGet {
		GetDataHandler(w, r)
	}

	incorrectRequestResponse(w, r)
}

// Handler return http.Handler to serve service requests
func Handler() http.Handler {
	apiMux := http.NewServeMux()
	return apiMux
}

func incorrectRequestResponse(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{ "error": "Invalid method"}`))
}
