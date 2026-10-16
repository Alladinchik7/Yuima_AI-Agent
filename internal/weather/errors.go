package weather

import "errors"

// Sentinel-ошибки для errors.Is на верхних слоях.
var (
	ErrCityNotFound   = errors.New("city not found")
	ErrUpstreamFailed = errors.New("weather upstream failed")
	ErrInvalidPayload = errors.New("invalid upstream payload")
)