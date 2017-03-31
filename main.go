package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/invite", invite)

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
		fmt.Fprintf(w, "This slack is only for women and gender minorities.")
		return
	}

	err = inviteUser(em)
	if err != nil {
		badRequest(w, err)
		return
	}

	fmt.Fprint(w, "Success! Check your email for invite.")
}

func badRequest(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "Bad request.", http.StatusBadRequest)
}
