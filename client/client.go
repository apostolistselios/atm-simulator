package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/apostolistselios/atm-simulator/api"
)

func main() {
	fmt.Println("WELCOME TO THE ATM.")
	username, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}

	// Make the user request to the server.
	if err := requestToUser(username); err != nil {
		log.Fatal(err)
	}

	answer := "y"
	for answer == "y" || answer == "Y" {
		userPrompt()
		var choice string
		if _, err := fmt.Scanf("%s\n", &choice); err != nil {
			log.Println("error incorrect choice")
			continue
		}

		choice = strings.ToLower(choice)
		if choice == "w" || choice == "d" {
			var amount int
			fmt.Print("PLEASE ENTER THE AMOUNT: ")
			if _, err := fmt.Scanf("%d\n", &amount); err != nil {
				log.Println("error incorrect amount")
				continue
			}

			// Check if the transaction is in the correct form.
			if err := checkTransaction(choice, amount); err != nil {
				log.Println(err)
				continue
			}

			// Make the transaction object.
			transaction := api.Transaction{
				UserID: username,
				Type:   choice,
				Amount: amount,
			}

			// Make the transaction request to the server.
			if err := requestToTransaction(transaction); err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("TRANSACTION COMPLETE")
		} else {
			// Make the balance request to the server.
			balance, err := requestToBalance(username)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("YOUR BALANCE IS:", balance)
		}
		// Check if the user wants to continue.
		fmt.Print("WOULD YOU LIKE TO CONTINUE (Y/N): ")
		fmt.Scanf("%s\n", &answer)
	}
}

// getCredentials gets the username from the user and returns it.
func getCredentials() (string, error) {
	var username string
	fmt.Print("USERNAME: ")
	if _, err := fmt.Scanf("%s\n", &username); err != nil {
		return "", errors.New("error parsing the username")
	}
	return username, nil
}

// requestToUser makes a HTTP request to /atm/user and handles the response.
func requestToUser(username string) error {
	url := "http://localhost:8080/atm/user"
	resp, err := http.Post(url, "text/plain", strings.NewReader(username))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	return nil
}

func userPrompt() {
	fmt.Println("1. W TO WITHDRAW AN AMOUNT")
	fmt.Println("2. D TO DEPOSIT AN AMOUNT")
	fmt.Println("3. B TO SEE YOUR BALANCE")
	fmt.Print("PLEASE CHOOSE THE TYPE A TYPE OF TRANSACTION: ")
}

// checkTransaction checks if the transaction is in the correct form.
// The transaction type (tranType) has to be w/W or d/D.
// The amount has to be a multiple of 20 or 50.
func checkTransaction(tranType string, amount int) error {
	if !(tranType == "w" || tranType == "W" || tranType == "d" || tranType == "D") {
		return errors.New("the transaction has to be between w/W or d/D, try again")
	}

	if amount%20 != 0 && amount%50 != 0 {
		return errors.New("the amount has to be multiple of 20 or multiple 50, try again")
	}
	return nil
}

// requestToTransaction makes a HTTP request to /atm/user/transaction and handles the response.
func requestToTransaction(transaction api.Transaction) error {
	// Encode the transaction object.
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(transaction); err != nil {
		return err
	}

	// Make the server request.
	url := "http://localhost:8080/atm/user/transaction"
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		data, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	return nil
}

// requestToBalance makes a HTTP request to /atm/user/balance and handles the response.
func requestToBalance(username string) (string, error) {
	url := "http://localhost:8080/atm/user/balance"
	resp, err := http.Post(url, "text/plain", strings.NewReader(username))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		return "", errors.New(string(data))
	}
	return string(data), nil
}
