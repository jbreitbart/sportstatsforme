// +build appengine

package email

import (
	"bufio"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"
)

func Breaststroke(w http.ResponseWriter, r *http.Request, emailText string, k *datastore.Key) {
	c := appengine.NewContext(r)
	/*	if err != nil {
		c.Errorf("No @ in email to address??? %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}*/

	reader := strings.NewReader(emailText)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		line := scanner.Text()
		c.Infof("%v", line)

	}
}

func Crawl() {

}
