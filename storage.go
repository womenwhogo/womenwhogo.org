package main

import (
	"io"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Served-With-Love-By", "Gophers")

	bucket := client.Bucket(bucketName)

	path := r.URL.Path[1:] // removing the leading / here

	if path == "" {
		path = "index.html"
	}

	object := bucket.Object(path)

	attr, attrErr := object.Attrs(r.Context())
	if attrErr == nil {
		w.Header().Set("Content-Type", attr.ContentType)
	}

	rc, err := object.NewReader(r.Context())
	if err != nil {
		http.Error(w, "404 Not Found.", http.StatusNotFound) // we need a fancy 404 here
		return
	}
	defer rc.Close()
	io.Copy(w, rc)
}

func handleAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=86400")
	http.FileServer(http.Dir("./womenwhogo.org/")).ServeHTTP(w, r)
}
