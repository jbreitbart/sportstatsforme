// +build appengine

package data

import (
	"time"

	"appengine"
	"appengine/datastore"
)

// SwimStyle represents the swim style
type SwimStyle int

const (
	Breaststroke SwimStyle = iota
	Crawl
)

// SwimStats is all information stored for swimming
type SwimStats struct {
	Style        SwimStyle      `datastore:",index"`
	Lanes        int            `datastore:",noindex"`
	LaneLength   int            `datastore:",noindex"`
	Time         time.Duration  `datastore:",noindex"`
	CreatedAt    time.Time      `datastore:",index"`
	User         *datastore.Key `datastore:",index"`
	DatastoreKey *datastore.Key `datastore:"-"`
}

// Store stores the SwimStats in the database
func (ss *SwimStats) Store(c appengine.Context) error {

	k := datastore.NewIncompleteKey(c, "SwimStats", nil)
	k, err := datastore.Put(c, k, ss)

	if err != nil {
		c.Errorf("Error while storing swim stats in datastore. Stats: %v. Error: %v", ss, err)
		ss.DatastoreKey = nil
	}

	ss.DatastoreKey = k

	return nil
}

// Delete removes a SwimStat from the database
func (ss *SwimStats) Delete(c appengine.Context) error {
	k := ss.DatastoreKey

	err := datastore.Delete(c, k)
	if err != nil {
		c.Errorf("Error while deleting SwimStat from datastore. Error: %v", err)
	}

	return err
}
