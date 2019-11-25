package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	http.HandleFunc("/invite", invite)
	http.HandleFunc("/assets/", handleAssets)
	http.HandleFunc("/", handleStatic)

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func invite(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		badRequest(w, r, err)
		return
	}

	em := r.Form.Get("email")
	okstr := r.Form.Get("ok")

	// ok references whether the individual identifies as a woman
	// or a gender minority.
	ok, err := strconv.ParseBool(okstr)
	if err != nil {
		badRequest(w, r, err)
		return
	}

	if !ok {
		http.Redirect(w, r, "http://www.womenwhogo.org/failure.html", http.StatusFound)
		return
	}

	err = inviteUser(r, em)
	if err != nil {
		interalServerError(w, r, err)
		return
	}

	http.Redirect(w, r, "http://www.womenwhogo.org/success.html", http.StatusFound)
}

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Error bad request: %v", err)
	http.Error(w, "Bad request.", http.StatusBadRequest)
}

func interalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %v", err)
	http.Error(w, "Internal server error.", http.StatusInternalServerError)
}
