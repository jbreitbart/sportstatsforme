// +build appengine

package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jbreitbart/sportstatsforme/data"
	"github.com/jbreitbart/sportstatsforme/targets"

	"appengine"
)

// Dispatch executes all commands for the user target
func Dispatch(w http.ResponseWriter, r *http.Request, u *data.User, urlPart *string) {
	c := appengine.NewContext(r)

	// get command
	command := strings.ToLower(targets.GetToken(urlPart))
	if command == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if command == "delete" {
		c.Infof("user/delete")

		u.SecureKey = targets.GetToken(urlPart)

		if err := u.Delete(c); err != nil {
			c.Errorf("Error at in user/delete. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User deleted!")

		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
