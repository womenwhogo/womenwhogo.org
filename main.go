package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const form = `
<html>
<body>
<form method="POST" action="/invite">
<label for="email">Email</label>
<input type="email" name="email" required></input>
<br>
<label for="ok">Do you identify as a woman or gender minority?</label>
<input type="radio" name="ok" value="true" required>Yes</input>
<input type="radio" name="ok" value="false" required>No</input>
<br>
<input type="submit" value="submit"></input>
</form>
</body>
</html>
`

func main() {
	http.HandleFunc("/invite", invite)
	http.HandleFunc("/", index)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	_, err := fmt.Fprint(w, form)
	if err != nil {
		log.Println(err)
	}
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
