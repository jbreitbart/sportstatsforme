// +build appengine

package main

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/jbreitbart/sportstatsforme/targets/email"
	"github.com/jhillyerd/go.enmime"

	"appengine"
)

func handlerEmails(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		c.Errorf("Error reading email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mime, err := enmime.ParseMIMEBody(msg)
	if err != nil {
		c.Errorf("Error parsing email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	from, err := msg.Header.AddressList("From")
	if err != nil {
		c.Errorf("Error reading email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// search for from address in DB
	u, err := getUser(w, r, from[0].Address)

	if err != nil {
		return
	}

	// u == user
	c.Infof("%v", u)

	to, err := msg.Header.AddressList("To")
	if err != nil {
		c.Errorf("Error reading email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emailAddress := to[0].Address
	at := strings.Index(emailAddress, "@")
	if at == -1 {
		c.Errorf("No @ in email to address??? %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO email in kleinbuchstaben umwandeln
	switch emailAddress[:at] {
	case "breaststroke":
		email.Breaststroke(w, r, mime.Text, u)
	case "brustschwimmen":
		email.Breaststroke(w, r, mime.Text, u)
	case "crawl":
		email.Crawl()
	}

	// parse email
	// string@appid.appspotmail.com

	// {breaststroke|crawl}@...
	// <lanes>*<lanesize>
	// time
}
