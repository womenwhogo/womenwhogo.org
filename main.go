package main

import (
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/invite", invite)
	http.HandleFunc("/assets/", handleAssets)
	http.HandleFunc("/", handleStatic)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func invite(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		badRequest(w, err)
		return
	}

	em := r.Form.Get("email")
	okstr := r.Form.Get("ok")

	// ok references whether the individual identifies as a woman
	// or a gender minority.
	ok, err := strconv.ParseBool(okstr)
	if err != nil {
		badRequest(w, err)
		return
	}

	if !ok {
		http.Redirect(w, r, "http://www.womenwhogo.org/failure.html", http.StatusFound)
		return
	}

	err = inviteUser(em)
	if err != nil {
		badRequest(w, err)
		return
	}

	http.Redirect(w, r, "http://www.womenwhogo.org/success.html", http.StatusFound)
}

func badRequest(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "Bad request.", http.StatusBadRequest)
}
