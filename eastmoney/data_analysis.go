package eastmoney

import "time"

// 数据分析

//计算某只基金，从from 买入，to时卖出 的收益率
// 期间会考虑这段时间的分红 (现金分红以红利再投资方式处理，份额折算则算为份额继续持有)
// 假如from 当天是假期，则按假期后一日来计算(天天基金网上的收益按假期前一日计算的)
func (fd Fund) CalcFundYield(from, to time.Time) float64 {
	yield, _ := fd.calcFundYield(from, to)
	return yield
}

// 返回收益率,及 实际开始投资日,
// 比如 基金于2016-10-1月为成立日, 而参数from日期可能早于此值如2016-9-1,则实际计算的时候按from=2016-10-1来计算,
// 此时startFundValue即为2016-10-1
func (fd Fund) calcFundYield(from, to time.Time) (yield float64, startFundValue FundValue) {
	var baseValue float64   // 起投那天的净值
	var cnt float64         // 持有份额
	var inMoney float64 = 1 // 按投入一玩计算
	var outMoney float64
	if from.Year() == to.Year() && from.Month() == to.Month() && from.Day() == to.Day() { // from==to
		return 0, startFundValue
	}
	if from.After(to) {
		return 0, startFundValue
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
			startFundValue = fv
			continue
		}
		if fv.FenHongType == FenHongType1 { // 1.每份基金份额折算1.012175663份 (折算之后 用户持有份额会增加，净值相应减少)
			cnt *= fv.FenHongRatio // (因折算导致 )持有份额增加
		} else if fv.FenHongType == FenHongType2 { // 2.每份派现金0.2150元,
			if fv.Value == 0 { // 有时返回的数据会是空行,则派现后,净值可能会变成1 2016-4-21 日的记录 如http://fund.eastmoney.com/f10/jjjz_001368.html
				fv.Value = 1
			}

			cnt += (cnt * fv.FenHongRatio) / fv.Value
		} else if fv.FenHongType == FenHongType3 { // 3. 每份基金份额分拆1.162668813份
			cnt *= fv.FenHongRatio // 持有份额增加
		}
		outMoney = cnt * fv.Value
	}
	yield = (outMoney - inMoney) / inMoney
	return
}

// 计算近1月收益率
func (fd Fund) CalcFundYieldLastMonth() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 30)
	return fd.CalcFundYield(from, to)
}

// 计算近2月收益率
func (fd Fund) CalcFundYieldLast2Month() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 60)
	return fd.CalcFundYield(from, to)
}

// 计算近3月收益率
func (fd Fund) CalcFundYieldLast3Month() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 90)
	return fd.CalcFundYield(from, to)
}

// 计算近6月收益率
func (fd Fund) CalcFundYieldLast6Month() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 180)
	return fd.CalcFundYield(from, to)
}

// 计算近1年收益率
func (fd Fund) CalcFundYieldLastYear() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365)
	return fd.CalcFundYield(from, to)
}

// 计算近2年收益率
func (fd Fund) CalcFundYieldLast2Year() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 2)
	return fd.CalcFundYield(from, to)
}

// 计算近3年收益率
func (fd Fund) CalcFundYieldLast3Year() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 3)
	return fd.CalcFundYield(from, to)
}

// 计算近5年收益率
func (fd Fund) CalcFundYieldLast5Year() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 5)
	return fd.CalcFundYield(from, to)
}

// 计算近10年收益率
func (fd Fund) CalcFundYieldLast10Year() float64 {
	now := time.Now()
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	from := to.Add(-time.Minute * 60 * 24 * 365 * 10)
	return fd.CalcFundYield(from, to)
}
