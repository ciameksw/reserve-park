package user

import (
	"net/http"

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

func (us *UserService) Login(r *http.Request) (*http.Response, error) {
	// Construct the full URL for the user service
	userServiceURL := us.UserURL + "/users/login"

	// Create a new HTTP request to the user service
	req, err := http.NewRequest(r.Method, userServiceURL, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	// Send the request to the user service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
