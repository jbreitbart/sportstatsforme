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
	Name         string `datastore:",noindex"`
	SecureKey    string `datastore:",noindex"`
	EMailAddress string `datastore:",index"`
}

// StoreUser stores the user in the database
func (u *User) Store(c appengine.Context) (*datastore.Key, error) {

	u.SecureKey = generateSecureKey()

	k := datastore.NewKey(c, "User", "", 0, nil)
	k, err := datastore.Put(c, k, u)

	if err != nil {
		c.Errorf("Error while storing user in datastore. User: %v. Error: %v", u, err)
	}

	return k, err
}

// DeleteUser removes a user and all its data from the database
func DeleteUser(c appengine.Context, k *datastore.Key, secureKey string) error {

	u := GetUserByKey(c, k)
	if u == nil {
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

// GetUserByKey returns the user based on its key
func GetUserByKey(c appengine.Context, key *datastore.Key) *User {
	var u User
	if err := datastore.Get(c, key, &u); err != nil {
		c.Errorf("Error while getting user based on key: %v", err)
		return nil
	}

	return &u
}

// Used to initialize the seed for random just once
var randomSeedInit sync.Once

func generateSecureKey() string {
	randomSeedInit.Do(func() {
		rand.Seed(time.Now().UTC().UnixNano())
	})

	keySize := 32
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, keySize)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
