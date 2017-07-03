package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/bundle.js", bundle)

	http.ListenAndServe(":9999", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("IIII")
	http.ServeFile(w, r, "index.html")
}

func bundle(w http.ResponseWriter, r *http.Request) {
	log.Println("AAAA")
	http.ServeFile(w, r, "dist/bundle.js")
}
