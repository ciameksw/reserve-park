package main

import (
	"github.com/ciameksw/reserve-park/facade/internal/facade/config"
	"github.com/ciameksw/reserve-park/facade/internal/facade/logger"
	"github.com/ciameksw/reserve-park/facade/internal/facade/server"
	"github.com/ciameksw/reserve-park/facade/internal/facade/services/reservation"
	"github.com/ciameksw/reserve-park/facade/internal/facade/services/spot"
	"github.com/ciameksw/reserve-park/facade/internal/facade/services/user"
)

func main() {
	// Get logger
	lgr := logger.GetLogger()

	// Get config
	cfg := config.GetConfig()

	// Get user service
	usr := user.NewUserService(cfg)

	// Get spot service
	spt := spot.NewSpotService(cfg)

	// Get reservation service
	rsrv := reservation.NewReservationService(cfg)

	s := server.NewServer(lgr, cfg, usr, spt, rsrv)
	s.Start()
}
