package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/apostolistselios/atm-simulator/api"
)

var store *api.Store

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error opening the server %s", err)
	}
	defer listener.Close()

	store, err = api.GetStore("../database.db")
	if err != nil {
		log.Fatalf("error connecting to database %s", err)
	}
	defer store.Close()

	fmt.Println("Server is listening on port:8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		handleConnection(conn)
	}
}

// handleConnection handles a client once the connection with the server is established
func handleConnection(conn net.Conn) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
	}

	_, err = checkUsername(strings.Trim(msg, "\n"))
	if err != nil {
		fmt.Fprintln(conn, err)
	}

	fmt.Fprintln(conn, "successful log in")
}

func checkUsername(username string) (api.User, error) {
	u, err := store.GetUser(username)
	return u, err
}

// func main() {
// 	store, err := api.GetStore("../database.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = store.GetUser("u")
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	store.ViewUsers()
// 	err = store.UpdateBalance("user4", 2000, "deposit")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	store.ViewUsers()
// }
