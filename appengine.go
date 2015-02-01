// +build appengine

package main

import "net/http"

func init() {
	http.HandleFunc("/u/", handlerWWW)
	http.HandleFunc("/_ah/mail/", handlerEmails)
	http.HandleFunc("/_ah/warmup", handlerWarmup)
}
