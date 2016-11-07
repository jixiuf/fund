package eastmoney

import (
	"time"

	"bitbucket.org/jixiuf/fund/utils"
)

// 定投基金收益率

type Period int

const (
	Week  Period = 1 // 每周定投一次
	Week2 Period = 2 // 每2周定投一次
	Month Period = 3 // 每月
)

// 数据分析
//计算某只基金，从from 开始定投，to时卖出 的收益率
// period 为定投周期,如果某定投日为节假日,则顺延到下一天定投
// 若按周定投,若起投日为周六日,则相当于起投日为下周一
func (fd Fund) CalcFundPeroidYield(period Period, from, to time.Time) (yield float64) {
	var cnt float64 = 0    // 定投次数
	var totalMoney float64 // 定投结束后 本金加利息
	if from.After(to) {
		return
	}
	var day time.Time = from
	for {
		if day.After(to) || (day.Year() == to.Year() && day.Month() == to.Month() && day.Day() == to.Day()) {
			break
		}
		if isPeriodDay(period, day, from, to) { // 如果day当天是定投日
			tmpYield := fd.CalcFundYield(day, to) // 计算day时买入1玩,到to日的收益率
			totalMoney += 1 * (1 + tmpYield)
			cnt++
		}
		day = day.Add(time.Minute * 60 * 24)
	}
	// cnt*1 为定期期间投入的本金总合
	// 而totalMoney 为 本金+利息
	yield = (totalMoney - cnt*1) / (cnt * 1)

	return
}

// 从from 那天开始,以period 为定投周期,判断day 是否是定投日
func isPeriodDay(period Period, day, from, to time.Time) bool {

	if period == Week {
		if day.Weekday() == from.Weekday() {
			return true
		}
		return false
	}
	if period == Week2 { // 每两周定投一次
		if day.Weekday() == from.Weekday() && int(day.Sub(from).Hours()/24/7)%2 == 0 {
			return true
		}
		return false
	}
	if period == Month { // 每月定投
		if day.Day() == from.Day() {
			return true
		}
		dayMonthDayCnt := utils.GetMonthDayCount(day.Year(), int(day.Month())) // from那月最多有多少天
		if dayMonthDayCnt < from.Day() && day.Day() == dayMonthDayCnt {
			// 如果day 为本月最后一天, 且起投日对应的那在 在本月无对应日,则以本月最后一天为定投日
			// 比如,定投日为31号, 而2月4月等根本没有31号,则以当月最后一天为定投日
			return true
		}

		return false
	}
	return false
}
