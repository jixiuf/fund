package eastmoney

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsPeriodDay_Week(t *testing.T) {
	from := time.Date(2016, 8, 4, 0, 0, 0, 0, time.Local) // 周四
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)

	day := time.Date(2016, 8, 11, 0, 0, 0, 0, time.Local) // 周四
	assert.True(t, isPeriodDay(Week, day, from, to))
}

func TestIsPeriodDay_Week2(t *testing.T) {
	from := time.Date(2016, 8, 4, 0, 0, 0, 0, time.Local) // 周四
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)

	day := time.Date(2016, 8, 11, 0, 0, 0, 0, time.Local)  // 周四
	day2 := time.Date(2016, 8, 18, 0, 0, 0, 0, time.Local) // 周四
	assert.False(t, isPeriodDay(Week2, day, from, to))
	assert.True(t, isPeriodDay(Week2, day2, from, to))
}

func TestIsPeriodDay_Month(t *testing.T) {
	from := time.Date(2016, 1, 31, 0, 0, 0, 0, time.Local) // 31号起投
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)

	day := time.Date(2016, 2, 29, 0, 0, 0, 0, time.Local) //
	assert.True(t, isPeriodDay(Month, day, from, to))
	day2 := time.Date(2016, 4, 30, 0, 0, 0, 0, time.Local) //
	assert.True(t, isPeriodDay(Month, day2, from, to))
}
func TestIsPeriodDay_Month2(t *testing.T) {
	from := time.Date(2016, 1, 15, 0, 0, 0, 0, time.Local) //
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)

	day := time.Date(2016, 2, 15, 0, 0, 0, 0, time.Local) //
	assert.True(t, isPeriodDay(Month, day, from, to))
}

// 下面的代友主要用于验证定投收益计算 与天天基金网上的计算差异
// 因为不知的天天基金网的计算收益是按哪天定投,故,差异是存在的
// 比如金鹰稳健成长混合  2016-11-5日显示 近1年定投收益为23.81%
// 而我计算的为25.12
func TestCalcFundPeroidYield(t *testing.T) {
	from := time.Date(2015, 11, 4, 0, 0, 0, 0, time.Local) //
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)
	fd, _ := GetFund("210004", 0)
	monthYield, _ := fd.CalcFundPeroidYieldWithPeriod(Month, from, to)
	fmt.Println(monthYield)
}

func TestCalcFundPeroidYield2(t *testing.T) {
	from := time.Date(2014, 11, 4, 0, 0, 0, 0, time.Local) //
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)
	fd, _ := GetFund("210004", 0)
	monthYield, _ := fd.CalcFundPeroidYieldWithPeriod(Month, from, to)
	fmt.Println(monthYield)
}

func TestCalcFundPeroidYield3(t *testing.T) {
	from := time.Date(2013, 11, 4, 0, 0, 0, 0, time.Local) //
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)
	fd, _ := GetFund("210004", 0)
	monthYield, _ := fd.CalcFundPeroidYieldWithPeriod(Month, from, to)
	fmt.Println(monthYield)
}

func TestCalcFundPeroidYield001368(t *testing.T) {
	fd, _ := GetFund("360005", 0)
	monthYield, _ := fd.CalcFundPeriodYieldLastYear(Month)
	fmt.Println(monthYield)
}
