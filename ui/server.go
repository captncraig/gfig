package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bundle.js", bundle)
	http.HandleFunc("/api/metadata", metadata)
	http.ListenAndServe(":8001", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func bundle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "dist/bundle.js")
}

func metadata(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "examplesettings.json")
}
