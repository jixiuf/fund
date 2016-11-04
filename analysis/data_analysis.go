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
