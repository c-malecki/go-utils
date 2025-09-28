package float

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const precision = 8

var scale = math.Pow(10, float64(9))

// currently based on MySQL decimal where latitude = decimal(10,8) and longitude = decimal(11,8)
func FormatLatLongFromString(latStr, longStr string) (string, string, error) {
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return "", "", fmt.Errorf("invalid latitude: %v", err)
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return "", "", fmt.Errorf("invalid longitude: %v", err)
	}

	lat = math.Round(lat*scale) / scale
	long = math.Round(long*scale) / scale

	finalLat := padToPrecision(lat, precision)
	finalLong := padToPrecision(long, precision)

	return finalLat, finalLong, nil
}

// currently based on MySQL decimal where latitude = decimal(10,8) and longitude = decimal(11,8)
func FormatLatLongFromFloat(lat, long float64) (string, string, error) {
	if lat == 0 || long == 0 {
		return "", "", errors.New("lat or long is 0")
	}

	lat = math.Round(lat*scale) / scale
	long = math.Round(long*scale) / scale

	finalLat := padToPrecision(lat, precision)
	finalLong := padToPrecision(long, precision)

	return finalLat, finalLong, nil
}

func padToPrecision(value float64, precision int) string {
	parts := strings.Split(fmt.Sprintf("%.*f", precision, value), ".")
	if len(parts) == 2 {
		for len(parts[1]) < precision {
			parts[1] += "0"
		}
		return parts[0] + "." + parts[1]
	}
	return fmt.Sprintf("%.*f", precision, value)
}
