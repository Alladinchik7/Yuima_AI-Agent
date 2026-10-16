package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client — погодный клиент
type Client struct {
	httpClient *http.Client
}

// NewClient создаёт клиент с общим http.Client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// GetWeather запрашивает погоду.
func (c *Client) GetWeather(ctx context.Context, city string) (*WeatherData, error) {
	baseURL := "https://ru.wttr.in" + url.PathEscape(city)
	params := url.Values{}
	params.Set("format", "j1")
	params.Set("lang", "ru")

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("http request creation failed: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUpstreamFailed, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%w: status=%d body=%s", ErrUpstreamFailed, res.StatusCode, string(body))
	}

	var wttr wttrResponse
	if err := json.NewDecoder(res.Body).Decode(&wttr); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
	}

	if len(wttr.CurrentCondition) == 0 {
		return nil, fmt.Errorf("%w: %q", ErrCityNotFound, city)
	}

	cur := wttr.CurrentCondition[0]

	var condition string
	if len(cur.WeatherDesc) > 0 {
		condition = cur.WeatherDesc[0].Value
	}

	var lat, lon float64
	if len(wttr.NearestArea) > 0 {
		lat = wttr.NearestArea[0].Latitude
		lon = wttr.NearestArea[0].Longitude
	}

	weather := &WeatherData{}
	weather.Fact.Temp = cur.TempC
	weather.Fact.FeelsLike = cur.FeelsLikeC
	weather.Fact.Condition = condition
	weather.Fact.WindSpeed = cur.WindspeedKmph / 3.6 // Конвертация в м/с
	weather.Fact.Humidity = cur.Humidity
	weather.Info.Lat = lat
	weather.Info.Lon = lon

	return weather, nil
}

// geocode — приватный, наружу не торчит.
/*func (c *Client) geocode(ctx context.Context, location string) (*Location, error) {
	baseURL := "https://geocoding-api.open-meteo.com/v1/search"
	params := url.Values{}
	params.Set("name", location)
	params.Set("country", "1")
	params.Set("language", "ru")

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("geocode request creation failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocode request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocode request failed: status code %d", resp.StatusCode)
	}

	var geoResp GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return nil, fmt.Errorf("geocode decode failed: %w", err)
	}

	if len(geoResp.Results) == 0 {
		return nil, fmt.Errorf("%w: %q", ErrCityNotFound, location)
	}

	return &geoResp.Results[0], nil
}*/