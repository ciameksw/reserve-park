package main

import (
	"github.com/ciameksw/reserve-park/user/internal/user/config"
	"github.com/ciameksw/reserve-park/user/internal/user/logger"
	"github.com/ciameksw/reserve-park/user/internal/user/mongodb"
	"github.com/ciameksw/reserve-park/user/internal/user/server"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg.MongoURI, "users") //TODO think about db and collection
	if err != nil {
		lgr.Error.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	s := server.NewServer(lgr, cfg, db)
	s.Start()
}
