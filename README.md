# ATM Simulator

### Description

An HTTP Server - Client application that simulates the functionality of an ATM.
The database uses a boltDB key value store and contains some sample users to test the functionality.

The Client can execute some basic operations. First you have to log in as user1, user2, user3, or user4
when prompted with are some demo users to demostrate the functionality.
Then you can choose to Withdraw (W), Deposit (D) from the account or see your balance (B).

### Running the Server

In order to run the server you have to be in the `./atm-simulator/server/` directory and execute
the following command in your terminal:

`go run server.go`

### Running the Client

In order to run the client you have to be in the `./atm-simulator/client/` directory and execute
the following command in your terminal:

`go run client.go`
