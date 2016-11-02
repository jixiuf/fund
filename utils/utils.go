package utils

import "strconv"

func Str2Float64(str string, defaultValue float64) (v float64) {
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultValue
	}

	return v
}
