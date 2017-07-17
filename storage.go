package main

import (
	"io"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

var (
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
)

func init() {
	var err error
	client, err = storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	bucketName = os.Getenv("SITE_BUCKET")
	if bucketName == "" {
		panic("No SITE_BUCKET given")
	}
}

func handleStaticRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Served-With-Love-By", "Gophers")

	bucket := client.Bucket(bucketName)

	path := r.URL.Path[1:] // removing the leading / here

	if path == "" {
		path = "index.html"
	}

	if path[:7] == "assets/" {
		w.Header().Set("Cache-Control", "max-age=86400")
	}

	object := bucket.Object(path)

	attr, attrErr := object.Attrs(r.Context())
	if attrErr == nil {
		w.Header().Set("Content-Type", attr.ContentType)
	}

	rc, err := object.NewReader(r.Context())
	defer rc.Close()
	if err != nil {
		http.Error(w, "404 Not Found.", http.StatusNotFound) // we need a fancy 404 here
		return
	}
	io.Copy(w, rc)
}
