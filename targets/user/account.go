// +build appengine

package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jbreitbart/sportstatsforme/data"
	"github.com/jbreitbart/sportstatsforme/targets"

	"appengine"
	"appengine/datastore"
)

// Dispatch executes all commands for the user target
func Dispatch(w http.ResponseWriter, r *http.Request, userKey *datastore.Key, urlPart *string) {
	c := appengine.NewContext(r)

	// get command
	command := strings.ToLower(targets.GetToken(urlPart))
	if command == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if command == "delete" {
		c.Infof("account/delete")

		secureKey := targets.GetToken(urlPart)

		if err := data.DeleteUser(c, userKey, secureKey); err != nil {
			c.Errorf("Error at in account/delete. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User deleted!")

		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
