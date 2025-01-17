package main

import (
	"github.com/ciameksw/reserve-park/spot/internal/spot/config"
	"github.com/ciameksw/reserve-park/spot/internal/spot/logger"
	"github.com/ciameksw/reserve-park/spot/internal/spot/mongodb"
	"github.com/ciameksw/reserve-park/spot/internal/spot/server"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg.MongoURI, "spots")
	if err != nil {
		lgr.Error.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect()

	s := server.NewServer(lgr, cfg, db)
	s.Start()
}
