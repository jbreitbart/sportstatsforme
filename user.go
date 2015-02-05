// +build appengine

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appengine"
	"appengine/mail"
	"appengine/urlfetch"

	"github.com/jbreitbart/sportstatsforme/data"
)

const accountCreationMessage = `
We have created an account with the following data just for you:

----
Name: %s
E-Mail: %s
----

Click the following link to delete your account
%s
`

func getUser(w http.ResponseWriter, r *http.Request, emailAddress string) (*data.User, error) {
	c := appengine.NewContext(r)

	u := data.GetUserByEmail(c, emailAddress)

	if u != nil {
		return u, nil
	}
	// not found => create an entry

	u = new(data.User)
	u.EMailAddress = emailAddress

	var err error
	if u.Name, err = getRandomName(w, r); err != nil {
		return nil, err
	}

	if err := u.Store(c); err != nil {
		c.Errorf("Error at storing user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	url := "http://sportstatsforme.appspot.com/u/" + u.DatastoreKey.Encode() + "/user/delete/" + u.SecureKey
	msg := &mail.Message{
		Sender:  "Sport Stats for Me <dontcare@sportstatsforme.appspotmail.com>",
		To:      []string{emailAddress},
		Subject: "Account created",
		Body:    fmt.Sprintf(accountCreationMessage, u.Name, u.EMailAddress, url),
	}

	c.Infof("Body: %v", msg.Body)

	if err := mail.Send(c, msg); err != nil {
		c.Errorf("Couldn't send email: %v", err)
	}

	return u, nil
}

func getRandomName(w http.ResponseWriter, r *http.Request) (string, error) {
	c := appengine.NewContext(r)

	// we need to generate a name first
	// get the json from
	// http://namey.muffinlabs.com/name.json?count=1&with_surname=true&frequency=rare
	// hmm, muffins
	client := urlfetch.Client(c)
	resp, err := client.Get("http://namey.muffinlabs.com/name.json?count=1&with_surname=true&frequency=rare")
	if err != nil {
		c.Errorf("Problem with Name generating request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}

	dec := json.NewDecoder(resp.Body)

	var m []string
	if err := dec.Decode(&m); err != nil {
		c.Errorf("Problem at decoding json request for the name: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}

	return m[0], nil
}
