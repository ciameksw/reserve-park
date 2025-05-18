package reservation

import (
	"bytes"
	"net/http"

	"github.com/ciameksw/reserve-park/facade/internal/facade/config"
	"github.com/ciameksw/reserve-park/facade/internal/facade/httpclient"
)

type ReservationService struct {
	ReservationURL string
}

func NewReservationService(cfg *config.Config) *ReservationService {
	return &ReservationService{
		ReservationURL: cfg.ReservationURL,
	}
}

func (rs *ReservationService) GetAll(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         rs.ReservationURL + "/reservations",
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

func (rs *ReservationService) DeleteReservation(reservationID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    rs.ReservationURL + "/reservations/" + reservationID,
		Method: http.MethodDelete,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (rs *ReservationService) GetReservationsBySpot(spotID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    rs.ReservationURL + "/reservations/spot/" + spotID,
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (rs *ReservationService) GetReservationsByUser(userID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    rs.ReservationURL + "/reservations/user/" + userID,
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (rs *ReservationService) Get(reservationID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    rs.ReservationURL + "/reservations/" + reservationID,
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (rs *ReservationService) Add(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         rs.ReservationURL + "/reservations",
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

func (rs *ReservationService) Edit(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         rs.ReservationURL + "/reservations",
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
