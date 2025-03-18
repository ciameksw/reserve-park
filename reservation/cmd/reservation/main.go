package main

import (
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/config"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/logger"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/mongodb"
	"github.com/ciameksw/reserve-park/reservation/internal/reservation/server"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg.MongoURI, "reservations")
	if err != nil {
		lgr.Error.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	s := server.NewServer(lgr, cfg, db)
	s.Start()
}
