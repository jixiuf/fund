package analysis

import (
	"time"

	"bitbucket.org/jixiuf/fund/eastmoney"
)

// 数据分析

//计算某只基金，从from 买入，to时卖出 的收益率
// 期间会考虑这段时间的分红 (现金分红以红利再投资方式处理，份额折算则算为份额继续持有)
// 假如from 当天是假期，则按假期后一日来计算(天天基金网上的收益按假期前一日计算的)
func CalcFundYield(fd eastmoney.Fund, from, to time.Time) float64 {
	var baseValue float64   // 起投那天的净值
	var cnt float64         // 持有份额
	var inMoney float64 = 1 // 按投入一玩计算
	var outMoney float64
	if from.Year() == to.Year() && from.Month() == to.Month() && from.Day() == to.Day() { // from==to
		return 0
	}
	if from.After(to) {
		return 0
	}

	for _, fv := range fd.FundValueList {
		if fv.Time.Before(from) {
			continue
		}
		if fv.Time.After(to) { // 到期赎回
			break
		}
		if baseValue == 0 { // 买入的那天
			baseValue = fv.Value
			cnt = inMoney / baseValue
			continue
		}
		if fv.FenHongType == eastmoney.FenHongType1 { // 1.每份基金份额折算1.012175663份 (折算之后 用户持有份额会增加，净值相应减少)
			cnt *= fv.FenHongRatio // (因折算导致 )持有份额增加
		} else if fv.FenHongType == eastmoney.FenHongType2 { // 2.每份派现金0.2150元,
			cnt += fv.FenHongRatio / fv.Value
		} else if fv.FenHongType == eastmoney.FenHongType3 { // 3. 每份基金份额分拆1.162668813份
			cnt *= fv.FenHongRatio // 持有份额增加
		}
		outMoney = cnt * fv.Value
	}
	return (outMoney - inMoney) / inMoney

}

// 计算近1月收益率
func CalcFundYieldLastMonth(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 30)
	return CalcFundYield(fd, from, to)
}

// 计算近2月收益率
func CalcFundYieldLast2Month(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 60)
	return CalcFundYield(fd, from, to)
}

// 计算近3月收益率
func CalcFundYieldLast3Month(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 90)
	return CalcFundYield(fd, from, to)
}

// 计算近6月收益率
func CalcFundYieldLast6Month(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 180)
	return CalcFundYield(fd, from, to)
}

// 计算近1年收益率
func CalcFundYieldLastYear(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365)
	return CalcFundYield(fd, from, to)
}

// 计算近2年收益率
func CalcFundYieldLast2Year(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 2)
	return CalcFundYield(fd, from, to)
}

// 计算近3年收益率
func CalcFundYieldLast3Year(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 3)
	return CalcFundYield(fd, from, to)
}

// 计算近5年收益率
func CalcFundYieldLast5Year(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 5)
	return CalcFundYield(fd, from, to)
}

// 计算近10年收益率
func CalcFundYieldLast10Year(fd eastmoney.Fund) float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 10)
	return CalcFundYield(fd, from, to)
}
