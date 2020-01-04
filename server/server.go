package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/apostolistselios/atm-simulator/api"
	"github.com/gorilla/mux"
)

var database = api.GetDatabase("../database.db")

func main() {
	api.ViewUsers(database)
	// router := http.NewServeMux()
	router := mux.NewRouter()
	router.HandleFunc("/atm/{userID}", verifyUser).Methods("GET")
	router.HandleFunc("/atm/{userID}/balance", viewBalance).Methods("GET")
	router.HandleFunc("/atm/{userID}/transaction", executeTransaction).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

// validateUser handles requests on /atm/{id} and checks if the user
// received from the request exists in the database. If yes sends back "OK"
// if not sends an http req with status code 400 and the error that occured.
func verifyUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]
	if err := api.CheckUserID(database, strings.Trim(userID, "\n")); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "OK")
}

// viewBalance handles requests on /atm/{id}/balance.
// Gets the balance for the specific user from the database and sends it back.
func viewBalance(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]
	if err := api.CheckUserID(database, strings.Trim(userID, "\n")); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	balance, err := api.GetBalance(database, userID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, balance)
}

// executeTransaction handles requests on /atm/{id}/transaction.
// Executes the transaction that is sent by the user on the database.
// If the transaction succeeds sends back "OK", if not sends an http req
// with status code 400 and the error that occured.
func executeTransaction(w http.ResponseWriter, req *http.Request) {
	var transaction api.Transaction
	if err := json.NewDecoder(req.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := api.ExecuteTransaction(database, transaction); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprint(w, "OK")
	api.ViewUsers(database)
}
