package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/ciameksw/reserve-park/user/internal/user/server"
)

func main() {
	// Read command line arguments
	serverHost := flag.String("host", "localhost", "Server host")
	serverPort := flag.String("port", "3001", "Server port")
	mongoURI := flag.String("mURI", "mongodb://localhost:27017", "MongoDB URI")
	mongoDB := flag.String("mDB", "users", "MongoDB database")
	flag.Parse()

	// Connect to MongoDB
	mongodb.Connect(*mongoURI, *mongoDB)
	defer mongodb.Disconnect()

	// Start server
	s := server.GetServer(*serverHost, *serverPort)
	fmt.Printf("Server started at %s:%s\n", *serverHost, *serverPort)
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		fmt.Println("Failed to start server")
		log.Fatal(err)
	}
}
