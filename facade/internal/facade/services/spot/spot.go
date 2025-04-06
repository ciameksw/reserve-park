package spot

import (
	"net/http"

	"github.com/ciameksw/reserve-park/facade/internal/facade/config"
	"github.com/ciameksw/reserve-park/facade/internal/facade/httpclient"
)

type SpotService struct {
	SpotURL string
}

func NewSpotService(cfg *config.Config) *SpotService {
	return &SpotService{
		SpotURL: cfg.SpotURL,
	}
}

func (ss *SpotService) GetSpotPrice(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         ss.SpotURL + "/spots/price",
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

func (ss *SpotService) GetAll(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         ss.SpotURL + "/spots",
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

func (ss *SpotService) GetSpot(spotID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    ss.SpotURL + "/spots/" + spotID,
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ss *SpotService) DeleteSpot(spotID string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    ss.SpotURL + "/spots/" + spotID,
		Method: http.MethodDelete,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ss *SpotService) AddSpot(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         ss.SpotURL + "/spots",
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

func (ss *SpotService) EditSpot(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         ss.SpotURL + "/spots",
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
