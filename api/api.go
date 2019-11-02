package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

// User struct
type User struct {
	Balance int
}

// GetDatabase returns a Database object and error.
func GetDatabase(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	return db, err
}

// CloseDatabase closes the database connection.
func CloseDatabase(db *bolt.DB) {
	db.Close()
}

// InsertUser inserts a user into the boltDB.
func InsertUser(db *bolt.DB, username string, balance int) error {
	u := User{balance}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		encodedUser, err := json.Marshal(u)
		if err != nil {
			return fmt.Errorf("error encoding User: %s", username)
		}
		return bucket.Put([]byte(username), encodedUser)
	})
}

// UpdateBalance updates the balance of the user with that specific username.
// transaction slice [username, tranType, amount]
func UpdateBalance(db *bolt.DB, transaction []string) error {
	username := transaction[0]
	tranType := transaction[1]
	amount, err := strconv.Atoi(transaction[2])
	if err != nil {
		return errors.New("error parsing the amount from the transaction")
	}

	return db.Update(func(tx *bolt.Tx) error {
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

		if tranType == "w" {
			if u.Balance > 0 && u.Balance-amount >= 0 {
				u.Balance -= amount
			} else {
				return fmt.Errorf("error withdrawal amount too high, User: %s Balance: %d", username, u.Balance)
			}
		} else if tranType == "d" {
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
func GetUser(db *bolt.DB, username string) (User, error) {
	var u User
	err := db.View(func(tx *bolt.Tx) error {
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
func ViewUsers(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))

		bucket.ForEach(func(key, value []byte) error {
			fmt.Printf("key=%s, value=%s\n", key, value)
			return nil
		})
		return nil
	})
}
