// +build appengine

package main

import (
	"net/http"
	"strings"

	"github.com/jbreitbart/sportstatsforme/data"
	"github.com/jbreitbart/sportstatsforme/targets"
	"github.com/jbreitbart/sportstatsforme/targets/user"

	"appengine"
	"appengine/datastore"
)

func handlerWWW(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	url := r.URL.Path

	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	// remove '/u/'
	url = url[3:]

	// extract user
	userKeyString := targets.GetToken(&url)
	if userKeyString == "" {
		// no user redirect to /
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
		return
	}

	c.Infof("UserKeyString: %s", url)
	userKey, err := datastore.DecodeKey(userKeyString)
	if err != nil {
		c.Errorf("Error at decoding key: %s, %s", userKeyString, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var u data.User
	u.DatastoreKey = userKey
	if err := u.GetByKey(c); err != nil {
		c.Errorf("Invalid key: %s, %s", userKeyString, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// get target
	target := strings.ToLower(targets.GetToken(&url))
	/*if target == "" {
		w.Header().Set("Location", r.URL.String()+"/show/www")
		w.WriteHeader(http.StatusFound)
		return
	}*/

	switch target {
	case "user":
		user.Dispatch(w, r, &u, &url)
	}

}

//<domain>/u/<userkey>/user/delete <- delete the account
//<domain>/u/<userkey>/swim/json <- get json swim data
