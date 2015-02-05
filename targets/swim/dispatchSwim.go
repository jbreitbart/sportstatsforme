// +build appengine

package swim

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jbreitbart/sportstatsforme/data"
	"github.com/jbreitbart/sportstatsforme/targets"

	"appengine"
	"appengine/datastore"
)

// Dispatch executes all commands for the swim target
func Dispatch(w http.ResponseWriter, r *http.Request, u *data.User, urlPart *string) {
	c := appengine.NewContext(r)

	// get command
	command := strings.ToLower(targets.GetToken(urlPart))
	if command == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if command == "delete" {
		c.Infof("swim/delete")

		var ss data.SwimStats
		var err error

		ss.DatastoreKey, err = datastore.DecodeKey(targets.GetToken(urlPart))
		if err != nil {
			c.Errorf("Error while decoding datatstore key. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = ss.Delete(c)
		if err != nil {
			c.Errorf("Error at in swim/delete. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Swim stat deleted!")

		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
