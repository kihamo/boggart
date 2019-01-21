package owntracks

import (
	"math"
)

const (
	earthRadiusKm = 6371
)

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1 = degreesToRadians(lat1)
	lon1 = degreesToRadians(lon1)
	lat2 = degreesToRadians(lat2)
	lon2 = degreesToRadians(lon2)

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(deltaLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthRadiusKm * 1000
}
