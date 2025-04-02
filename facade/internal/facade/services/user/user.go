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
	userServiceURL := us.UserURL + "/users/login"

	req, err := http.NewRequest(r.Method, userServiceURL, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) GetAll(r *http.Request) (*http.Response, error) {
	userServiceURL := us.UserURL + "/users"

	req, err := http.NewRequest(r.Method, userServiceURL, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Authorize(authHeader string) (*http.Response, error) {
	userServiceURL := us.UserURL + "/users/authorize"

	req, err := http.NewRequest(http.MethodGet, userServiceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
