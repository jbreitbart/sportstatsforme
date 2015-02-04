// +build appengine

package data

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"appengine"
	"appengine/datastore"
)

// User is the datastore information for an individual user
type User struct {
	Name         string         `datastore:",noindex"`
	SecureKey    string         `datastore:",noindex"`
	EMailAddress string         `datastore:",index"`
	DatastoreKey *datastore.Key `datastore:"-"`
}

// Store stores the user in the database
func (u *User) Store(c appengine.Context) error {
	u.SecureKey = generateSecureKey()

	k := datastore.NewIncompleteKey(c, "User", nil)
	k, err := datastore.Put(c, k, u)

	if err != nil {
		c.Errorf("Error while storing user in datastore. User: %v. Error: %v", u, err)
		u.DatastoreKey = nil
	} else {
		u.DatastoreKey = k
	}

	return err
}

// Delete removes a user and all its data from the database
func (u *User) Delete(c appengine.Context) error {
	k := u.DatastoreKey
	secureKey := u.SecureKey

	if err := u.GetByKey(c); err != nil {
		c.Infof("User not found. Key invalid")
		return errors.New("Invalid key.")
	}

	// validate secure key
	if u.SecureKey != secureKey {
		c.Infof("Wrong secure key for deleting User %v", u.EMailAddress)
		return errors.New("Invalid key. " + secureKey + " - " + u.SecureKey)
	}

	err := datastore.Delete(c, k)
	if err != nil {
		c.Errorf("Error while deleting user from datastore. Error: %v", err)
	}

	// TODO delete every user data

	return err
}

// GetUserByEmail returns the user key based on the email
func GetUserByEmail(c appengine.Context, email string) *User {
	c.Infof("Searching for user: %v", email)

	q := datastore.NewQuery("User").
		Filter("EMailAddress =", email).
		Limit(1)

	var u []User
	keys, err := q.GetAll(c, &u)

	if err != nil {
		c.Errorf("Error at user key query: %v", err)
		return nil
	}

	if len(keys) == 0 {
		c.Infof("Not found.")
		return nil
	}

	returnee := u[0]
	return &returnee
}

// GetByKey returns the user based on its key
func (u *User) GetByKey(c appengine.Context) error {

	key := u.DatastoreKey

	if err := datastore.Get(c, key, u); err != nil {
		c.Errorf("Error while getting user based on key: %v", err)
		return err
	}

	return nil
}

// Used to initialize the seed for random just once
var randomSeedInit sync.Once

// Alphabet for the secure key
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Length of the secure key
const keySize = 32

func generateSecureKey() string {
	randomSeedInit.Do(func() {
		rand.Seed(time.Now().UTC().UnixNano())
	})

	b := make([]rune, keySize)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

/*
// GetUserKeyByEmail returns the user key based on the email
func GetUserKeyByEmail(c appengine.Context, email string) *datastore.Key {
	c.Infof("Searching for user: %v", email)

	q := datastore.NewQuery("User").
		Filter("EMailAddress =", email).
		Limit(1).
		KeysOnly()

	keys, err := q.GetAll(c, nil)
	if err != nil {
		c.Errorf("Error at user key query: %v", err)
		return nil
	}

	if len(keys) == 0 {
		c.Infof("Not found.")
		return nil
	}

	return keys[0]
}
*/
