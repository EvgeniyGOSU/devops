package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		Humidity int     `json:"humidity"`
		WindKph  float64 `json:"wind_kph"`
	} `json:"current"`
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetWeather(city string) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", c.baseURL, c.apiKey, city)

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус: %s", resp.Status)
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("ошибка парсинга: %w", err)
	}

	return &weather, nil
}
