// +build appengine

package data

import (
	"time"

	"appengine/datastore"
)

// SwimStyle represents the swim style
type SwimStyle int

const (
	breaststroke SwimStyle = iota
	crawl
)

// SwimStats is all information stored for swimming
type SwimStats struct {
	Style     SwimStyle     `datastore:",index"`
	Lanes     int           `datastore:",noindex"`
	LaneSize  int           `datastore:",noindex"`
	Time      time.Duration `datastore:",noindex"`
	CreatedAt time.Time     `datastore:",index"`
	User      datastore.Key `datastore:",index"`
}
