// +build appengine

package email

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jbreitbart/sportstatsforme/data"

	"appengine"
	"appengine/mail"
)

const confirmationMessage = `
We have added the following statistics to your account:

----
Swim style: %v
#Lanes: %v
Lane length: %v
Time: %v
----

Click the following link to delete the stats
%s
`

// Breaststroke parses an email and stores the data in the datastore
// Format:
// <lane> * <lane length>
// time
func Breaststroke(w http.ResponseWriter, r *http.Request, emailText string, u *data.User) {
	c := appengine.NewContext(r)

	c.Infof("Will try to parse breakstroke email!")

	var ss data.SwimStats
	ss.Style = data.Breaststroke
	ss.CreatedAt = time.Now()
	ss.User = u.DatastoreKey

	/*      if err != nil {
	        c.Errorf("No @ in email to address??? %v", err)
	        w.WriteHeader(http.StatusInternalServerError)
	        return
	}*/

	reader := strings.NewReader(emailText)
	scanner := bufio.NewScanner(reader)

	count := 0
	for scanner.Scan() {

		line := scanner.Text()
		c.Infof("%v: %v", count, line)

		switch count {
		case 0:
			var lanesString, laneLengthString string
			if index := strings.Index(line, "*"); index != -1 {
				lanesString = line[:index]
				laneLengthString = line[index+1:]
			} else {
				lanesString = line
				laneLengthString = "50"
			}

			if number, err := strconv.ParseInt(lanesString, 0, 32); err != nil {
				c.Errorf("Could not parse number of lanes %v - %v", lanesString, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				ss.Lanes = int(number)
			}

			if number, err := strconv.ParseInt(laneLengthString, 0, 32); err != nil {
				c.Errorf("Could not parse length of lanes %v - %v", lanesString, err)
				ss.LaneLength = 50
			} else {
				ss.LaneLength = int(number)
			}

		case 1:
			if d, err := time.ParseDuration(line); err != nil {
				c.Infof("Could not parse duration: %v", err)
			} else {
				ss.Time = d
			}

		default:
			c.Infof("Ignoreing line #%v: %v", count, line)

		}
		count++
	}

	c.Infof("Swimstats: %v", ss)
	if err := ss.Store(c); err != nil {
		c.Errorf("Error while storing swimstat: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send confirmation email
	url := "http://sportstatsforme.appspot.com/u/" + u.DatastoreKey.Encode() + "/swim/delete/" + ss.DatastoreKey.Encode()
	msg := &mail.Message{
		Sender:  "Sport Stats for Me <dontcare@sportstatsforme.appspotmail.com>",
		To:      []string{u.EMailAddress},
		Subject: "Added stats to your user",
		Body:    fmt.Sprintf(confirmationMessage, "breaststroke", ss.Lanes, ss.LaneLength, ss.Time, url),
	}

	c.Infof("Body: %v", msg.Body)

	if err := mail.Send(c, msg); err != nil {
		c.Errorf("Couldn't send email: %v", err)
	}

}

func Crawl() {

}
