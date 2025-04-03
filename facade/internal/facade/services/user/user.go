package user

import (
	"bytes"
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

func (us *UserService) Register(body []byte) (*http.Response, error) {
	userServiceURL := us.UserURL + "/users"

	req, err := http.NewRequest("POST", userServiceURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func (us *UserService) GetUser(userID string) (*http.Response, error) {
	userServiceURL := us.UserURL + "/users/" + userID

	req, err := http.NewRequest("GET", userServiceURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) DeleteUser(userID string) (*http.Response, error) {
	userServiceURL := us.UserURL + "/users/" + userID

	req, err := http.NewRequest("DELETE", userServiceURL, nil)
	if err != nil {
		return nil, err
	}

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
