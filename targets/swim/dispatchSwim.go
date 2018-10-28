// +build appengine

package swim

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/jbreitbart/sportstatsforme/bindata"
	"github.com/jbreitbart/sportstatsforme/data"
	"github.com/jbreitbart/sportstatsforme/targets"

	"appengine"
	"appengine/datastore"
)

var statsAsset = bindata.MustAsset("swimstats.html")
var statsTemplate = template.Must(template.New("swimStats").Parse(string(statsAsset)))

// Dispatch executes all commands for the swim target
func Dispatch(w http.ResponseWriter, r *http.Request, u *data.User, urlPart *string) {

	c := appengine.NewContext(r)

	// get command
	command := strings.ToLower(targets.GetToken(urlPart))

	switch command {
	case "delete":
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

	case "stats":
		if err := statsTemplate.Execute(w, nil); err != nil {
			c.Errorf("Error rendering swim stats template. Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "json":
		command := strings.ToLower(targets.GetToken(urlPart))

		switch command {
			case "lanesPerWeek":
				s, err := lanesPerWeek(c, u)
				if err != nil {
					c.Errorf("Error at in lanesPerWeek. Error: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(w, s)

			default:
				w.WriteHeader(http.StatusNotImplemented)
		}


	default:
		w.WriteHeader(http.StatusNotImplemented)
	}

}

type lanesPerWeek struct {
	week int
	year int
	number int
}

func lanesPerWeek(c appengine.Context, u *data.User) (string, error) {
	ss, err := data.GetAllSwimStatsforUser(c, u)
	if err != nil {
		return "", err
	}


	for s := range ss {

	}
	return "", nil
}
