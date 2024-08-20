package main

import (
	"github.com/ciameksw/reserve-park/spot/internal/spot/config"
	"github.com/ciameksw/reserve-park/spot/internal/spot/logger"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()
}
