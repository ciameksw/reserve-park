package reservation

import "github.com/ciameksw/reserve-park/facade/internal/facade/config"

type ReservationService struct {
	ReservationURL string
}

func NewReservationService(cfg *config.Config) *ReservationService {
	return &ReservationService{
		ReservationURL: cfg.ReservationURL,
	}
}
