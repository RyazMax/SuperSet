package main

import (
	"log"
	"net/http"

	"../src/api"
	"../src/static"
	"../src/templates"
	"../src/universe"
)

func main() {
	err := universe.Init("127.0.0.1", 6666)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/", api.Handler())
	http.Handle("/static/", http.StripPrefix("/static/", static.Handler()))
	http.Handle("/", templates.Handler())

	log.Println("Server starting on 127.0.0.1:9999")
	log.Fatal(http.ListenAndServe("127.0.0.1:9999", nil))
}
