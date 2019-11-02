package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/apostolistselios/atm-simulator/api"
	"github.com/boltdb/bolt"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error opening the server %s", err)
	}
	defer listener.Close()

	database, err := api.GetDatabase("../database.db")
	if err != nil {
		log.Fatalf("error connecting to database %s", err)
	}
	defer api.CloseDatabase(database)

	fmt.Println("Server is listening on port:8080...")
	api.ViewUsers(database)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		handleConnection(conn, database)
	}
}

// handleConnection handles a client once the connection with the server is established
func handleConnection(conn net.Conn, database *bolt.DB) {
	// Receives the username from the client.
	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
	}

	// Check if the username is valid.
	_, err = api.GetUser(database, strings.Trim(username, "\n"))
	if err != nil {
		fmt.Fprintln(conn, err)
	}
	fmt.Fprintln(conn, "OK")

	// Receive transactions.
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			fmt.Fprintln(conn, err)
			continue
		}
		msg = strings.Trim(msg, "\n")

		if msg == "exit" {
			break
		}

		transaction := strings.Split(msg, " ")
		err = api.UpdateBalance(database, transaction)
		if err != nil {
			fmt.Fprintln(conn, err)
			continue
		}
		fmt.Fprintln(conn, "OK")
		api.ViewUsers(database)
	}
}
