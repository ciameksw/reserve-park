package user

import (
	"github.com/ciameksw/reserve-park/facade/internal/facade/config"
)

type UserService struct {
	UserURL string
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{
		UserURL: cfg.UserURL,
	}
}
