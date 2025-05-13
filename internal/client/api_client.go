package client

import (
	"encoding/json"
	"fmt"
	"music-library/internal/model"
	"net/http"
	"net/url"
)

type APIClient struct {
	BaseURL string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

func (c *APIClient) FetchSongDetails(group, song string) (*model.SongDetail, error) {
	urlS := fmt.Sprintf("%s/info?group=%s&song=%s",
		c.BaseURL,
		url.QueryEscape(group),
		url.QueryEscape(song),
	)

	resp, err := http.Get(urlS)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var detail model.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	return &detail, nil
}
