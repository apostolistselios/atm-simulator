package api

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

// User struct
type User struct {
	Balance int
}

// Store struct that holds a boldDB
type Store struct {
	db *bolt.DB
}

// GetStore returns a Store object and error.
func GetStore(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	return &Store{db: db}, err
}

// Close closes the database connection.
func (store *Store) Close() {
	store.db.Close()
}

// InsertUser inserts a user into the boltDB.
func (store *Store) InsertUser(username string, balance int) error {
	u := User{balance}

	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		encodedUser, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("error encoding User: %s", username)
		}
		return bucket.Put([]byte(username), encodedUser)
	})
}

// UpdateBalance updates the balance of the user with that specific username.
func (store *Store) UpdateBalance(username string, amount int, operation string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(username))
		if value == nil {
			return fmt.Errorf("error User: %s doesn't exist", username)
		}

		var u User
		err := json.Unmarshal(value, &u)
		if err != nil {
			return fmt.Errorf("error decoding User: %s", username)
		}

		if operation == "withdraw" {
			if u.Balance > 0 && u.Balance-amount >= 0 {
				u.Balance -= amount
			} else {
				return fmt.Errorf("error withdrawal amount too high, User: %s Balance: %d", username, u.Balance)
			}
		} else if operation == "deposit" {
			u.Balance += amount
		}

		encodedUser, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("error encoding User: %s", username)
		}
		return bucket.Put([]byte(username), encodedUser)
	})
}

// GetUser returns the User with that specific username.
func (store *Store) GetUser(username string) (User, error) {
	var u User
	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		value := bucket.Get([]byte(username))
		if value == nil {
			return fmt.Errorf("error User: %s doesn't exist", username)
		}

		err := json.Unmarshal(value, &u)
		if err != nil {
			return fmt.Errorf("error decoding User: %s", username)
		}
		return nil
	})

	return u, err
}

// ViewUsers prints all the users in the console.
func (store *Store) ViewUsers() {
	store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		bucket.ForEach(func(key, value []byte) error {
			fmt.Printf("key=%s, value=%s\n", key, value)
			return nil
		})
		return nil
	})
}
