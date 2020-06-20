package mat

import "math"

func Round(val float64) float64 {
	return math.Round(val*100) / 100
}

func RoundAcc(val float64, accuracy float64) float64 {
	acc := math.Pow(10, accuracy)
	return math.Round(val*acc) / acc
}
