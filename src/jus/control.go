package jus

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index."))
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("project."))
}
