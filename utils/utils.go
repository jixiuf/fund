package utils

import "strconv"

func Str2Float64(str string, defaultValue float64) (v float64) {
	v, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultValue
	}

	return v
}

var monthDays = [...]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// 获取指定年月一共有多少天
func GetMonthDayCount(year, month int) int { // month[1~12]
	if month == 2 && IsLeapYear(year) {
		return 29
	}
	return monthDays[month]
}
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}
