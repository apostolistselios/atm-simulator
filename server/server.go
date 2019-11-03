package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apostolistselios/atm-simulator/api"
)

var database = api.GetDatabase("../database.db")

func main() {
	api.ViewUsers(database)
	router := http.NewServeMux()
	router.HandleFunc("/atm/user", checkUserHandler)
	router.HandleFunc("/atm/user/balance", balanceHandler)
	router.HandleFunc("/atm/user/transaction", transactionHandler)
	http.ListenAndServe(":8080", router)
}

// checkUserHandler handles requests on /atm/user and checks if the user
// received from the request exists in the database. If yes sends back "OK"
// if not sends an http req with status code 400 and the error that occured.
func checkUserHandler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	username := string(body)
	err := checkUsername(username)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "OK")
}

// checkUsername checks if the specific username exists in the database.
func checkUsername(username string) error {
	_, err := api.GetUser(database, strings.Trim(username, "\n"))
	return err
}

// balanceHandler handles requests on /atm/user/balance.
// Gets the balance for the specific user from the database and sends it back.
func balanceHandler(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	username := string(body)
	balance, err := api.GetBalance(database, username)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, balance)
}

// transactionHandler handles requests on /atm/user/transaction.
// Executes the transaction that is sent by the user on the database.
// If the transaction succeeds sends back "OK", if not sends an http req
// with status code 400 and the error that occured.
func transactionHandler(w http.ResponseWriter, req *http.Request) {
	var transaction api.Transaction
	if req.Body == nil {
		http.Error(w, "error send a request body", 400)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = api.ExecuteTransaction(database, transaction)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "OK")
	api.ViewUsers(database)
}
