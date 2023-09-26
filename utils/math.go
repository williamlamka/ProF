package utils

import (
	"math"
)

func LimitDecimalDigits(value float32) float32 {
	return float32(math.Round(float64(value)*100) / 100)
}
