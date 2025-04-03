package user

import (
	"bytes"
	"io"
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
	ct := "application/json"
	params := sendRequestParams{
		Path:          "/users",
		Method:        http.MethodPost,
		Body:          bytes.NewBuffer(body),
		ContentType:   &ct,
		Authorization: nil,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Login(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := sendRequestParams{
		Path:          "/users/login",
		Method:        r.Method,
		Body:          r.Body,
		ContentType:   &ct,
		Authorization: nil,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Edit(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := sendRequestParams{
		Path:          "/users",
		Method:        http.MethodPatch,
		Body:          bytes.NewBuffer(body),
		ContentType:   &ct,
		Authorization: nil,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) GetAll(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := sendRequestParams{
		Path:          "/users",
		Method:        r.Method,
		Body:          r.Body,
		ContentType:   &ct,
		Authorization: nil,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) GetUser(userID string) (*http.Response, error) {
	params := sendRequestParams{
		Path:          "/users/" + userID,
		Method:        http.MethodGet,
		Body:          nil,
		ContentType:   nil,
		Authorization: nil,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) DeleteUser(userID string) (*http.Response, error) {
	params := sendRequestParams{
		Path:          "/users/" + userID,
		Method:        http.MethodDelete,
		Body:          nil,
		ContentType:   nil,
		Authorization: nil,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (us *UserService) Authorize(authHeader string) (*http.Response, error) {
	params := sendRequestParams{
		Path:          "/users/authorize",
		Method:        http.MethodGet,
		Body:          nil,
		ContentType:   nil,
		Authorization: &authHeader,
	}
	resp, err := us.sendRequestToUserService(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type sendRequestParams struct {
	Path          string
	Method        string
	Body          io.Reader
	ContentType   *string
	Authorization *string
}

func (us *UserService) sendRequestToUserService(params sendRequestParams) (*http.Response, error) {
	requestURL := us.UserURL + params.Path

	var body io.Reader
	if params.Body != nil {
		body = params.Body
	}
	req, err := http.NewRequest(params.Method, requestURL, body)
	if err != nil {
		return nil, err
	}

	if params.ContentType != nil {
		req.Header.Set("Content-Type", *params.ContentType)
	}
	if params.Authorization != nil {
		req.Header.Set("Authorization", *params.Authorization)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
