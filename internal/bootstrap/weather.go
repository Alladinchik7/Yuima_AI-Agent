package bootstrap

import "yuima-agent/internal/weather"

func InitWeatherClient() *weather.Client {
	return weather.NewClient()
}