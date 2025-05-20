package user

import (
	"bytes"
	"net/http"

	"github.com/ciameksw/reserve-park/facade/internal/facade/config"
	"github.com/ciameksw/reserve-park/facade/internal/facade/httpclient"
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
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         us.UserURL + "/users",
		Method:      http.MethodPost,
		Body:        bytes.NewBuffer(body),
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Login(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         us.UserURL + "/users/login",
		Method:      r.Method,
		Body:        r.Body,
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Edit(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         us.UserURL + "/users",
		Method:      http.MethodPatch,
		Body:        bytes.NewBuffer(body),
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) GetAll(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         us.UserURL + "/users",
		Method:      r.Method,
		Body:        r.Body,
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) GetUser(userID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    us.UserURL + "/users/" + userID,
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) DeleteUser(userID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    us.UserURL + "/users/" + userID,
		Method: http.MethodDelete,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Authorize(authHeader string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:           us.UserURL + "/users/authorize",
		Method:        http.MethodGet,
		Authorization: &authHeader,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
