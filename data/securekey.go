package data

import (
	"math/rand"
	"sync"
	"time"
)

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
