// +build appengine

package data

import (
	"errors"
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

// Returns all swimstats for a specific user
func GetAllSwimStatsforUser(c appengine.Context, u *User) ([]SwimStats, error) {
	c.Infof("Get all swim stats for User: %v", u)

	q := datastore.NewQuery("SwimStats").
		Filter("User =", u.DatastoreKey)

	var stats []SwimStats
	keys, err := q.GetAll(c, &stats)

	if err != nil {
		c.Errorf("Error at get all SwimStats for user: %v", err)
		return nil, err
	}

	if len(keys) == 0 {
		c.Infof("Found nothing.")
		return nil, nil
	}

	return stats, nil
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

func deleteAllSwimStatsofUser(c appengine.Context, userKey *datastore.Key) error {
	c.Infof("Delete all swimStats of a user.")
	q := datastore.NewQuery("SwimStats").
		Filter("User =", userKey).
		KeysOnly()

	keys, err := q.GetAll(c, nil)

	if err != nil {
		c.Errorf("Error while getting all swimStats keys of a user: %v", err)
		return err
	}

	if len(keys) == 0 {
		c.Infof("None found.")
		return errors.New("asd")
	}

	err = datastore.DeleteMulti(c, keys)
	if err != nil {
		c.Errorf("Error while deleting swimStats: %v", err)
		return err
	}

	return nil

}
