package main

import (
	"log"

	"github.com/ciameksw/reserve-park/user/internal/user/config"
	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/ciameksw/reserve-park/user/internal/user/server"
)

func main() {
	// Get config
	cfg := config.GetConfig()

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg.MongoURI, "users") //TODO think about db and collection
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	s := server.NewServer(cfg, db)
	s.Start()
}
