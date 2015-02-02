// +build appengine

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/mail"
	"appengine/urlfetch"

	"github.com/jbreitbart/sportstatsforme/data"
)

func userKey(w http.ResponseWriter, r *http.Request, emailAddress string) (*datastore.Key, error) {
	c := appengine.NewContext(r)
	k := data.GetUserKeyByEmail(c, emailAddress)

	if k != nil {
		return k, nil
	}

	// not found => create an entry

	// we need to generate a name first
	// get the json from
	// http://namey.muffinlabs.com/name.json?count=1&with_surname=true&frequency=rare
	// hmm, muffins
	client := urlfetch.Client(c)
	resp, err := client.Get("http://namey.muffinlabs.com/name.json?count=1&with_surname=true&frequency=rare")
	if err != nil {
		c.Errorf("Problem with Name generating request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	var u data.User
	var m []string
	if err := dec.Decode(&m); err != nil {
		c.Errorf("Problem at decoding json request for the name: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	// Store it
	u.Name = m[0]
	u.EMailAddress = emailAddress
	k, err = u.Store(c)
	if err != nil {
		c.Errorf("Error at storing user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	url := "http://sportstatsforme.appspot.com/u/" + k.Encode() + "/user/delete/" + u.SecureKey
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

	return k, nil
}

const accountCreationMessage = `
We have created an account with the following data just for you:

----
Name: %s
E-Mail: %s
----

Click the following link to delete your account
%s
`
