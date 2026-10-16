package weather

// WeatherData — ответ Яндекс.Погоды.
type WeatherData struct {
	
	Fact struct {
		Temp      int     `json:"temp"`
		FeelsLike int     `json:"feels_like"`
		Condition string  `json:"condition"`
		WindSpeed float64 `json:"wind_speed"`
		Humidity  int     `json:"humidity"`
	} `json:"fact"`
	Info struct {
		Lat float64 `json:"latitude"`
		Lon float64 `json:"longitude"`
	} `json:"info"`
}

// wttrResponse — ответ wttr.in в формате j1.
type wttrResponse struct {
	CurrentCondition []struct {
		TempC         int     `json:"temp_C,string"`
		FeelsLikeC    int     `json:"FeelsLikeC,string"`
		WeatherDesc   []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		WindspeedKmph float64 `json:"windspeedKmph,string"`
		Humidity      int     `json:"humidity,string"`
	} `json:"current_condition"`
	NearestArea []struct {
		Latitude  float64 `json:"latitude,string"`
		Longitude float64 `json:"longitude,string"`
	} `json:"nearest_area"`
}

/*  GeocodingResponse — ответ open-meteo geocoding.
type GeocodingResponse struct {
	Results []Location `json:"results"`
}

// Location — найденная локация.
type Location struct {
	Name     string  `json:"name"`
	Lat      float64 `json:"latitude"`
	Lon      float64 `json:"longitude"`
	Country  string  `json:"country"`
	Timezone string  `json:"timezone"`
	}*/