package main

import (
	"github.com/ciameksw/reserve-park/spot/internal/spot/config"
	"github.com/ciameksw/reserve-park/spot/internal/spot/logger"
)

func main() {
	// Get logger
	_ = logger.GetLogger()

	// Get config
	_ = config.GetConfig()

	// Connect to MongoDB
	// TODO
}
