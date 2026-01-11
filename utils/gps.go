package utils

import (
	"math"

	"glosindo-backend-go/config"
)

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // meters

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func IsInOfficeArea(lat, lng float64) bool {
	distance := CalculateDistance(lat, lng, config.AppConfig.OfficeLatitude, config.AppConfig.OfficeLongitude)
	return distance <= config.AppConfig.MaxDistanceMeters
}

func GetDistanceFromOffice(lat, lng float64) float64 {
	return CalculateDistance(lat, lng, config.AppConfig.OfficeLatitude, config.AppConfig.OfficeLongitude)
}