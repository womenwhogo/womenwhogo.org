package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const inviteURL = "https://slack.com/api/users.admin.invite"

var token = os.Getenv("SLACK_API_TOKEN")

func inviteUser(r *http.Request, em string) error {
	values := make(url.Values, 2)
	values.Add("token", token)
	values.Add("email", em)
	// resend invite if already sent
	values.Add("resend", "true")

	s := values.Encode()
	u := inviteURL + "?" + s

	resp, err := http.Get(u)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(b))
	}

	return nil
}
