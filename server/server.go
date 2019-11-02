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

func checkUserHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	username := string(body)
	err = checkUsername(username)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "OK")
}

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

	err = api.UpdateBalance(database, transaction)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "OK")
	api.ViewUsers(database)
}

func main() {
	api.ViewUsers(database)
	router := http.NewServeMux()
	router.HandleFunc("/atm/user", checkUserHandler)
	router.HandleFunc("/atm/user/transaction", transactionHandler)
	http.ListenAndServe(":8080", router)
}

func checkUsername(username string) error {
	_, err := api.GetUser(database, strings.Trim(username, "\n"))
	return err
}
