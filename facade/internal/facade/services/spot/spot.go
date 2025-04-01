package spot

import "github.com/ciameksw/reserve-park/facade/internal/facade/config"

type SpotService struct {
	SpotURL string
}

func NewSpotService(cfg *config.Config) *SpotService {
	return &SpotService{
		SpotURL: cfg.SpotURL,
	}
}
