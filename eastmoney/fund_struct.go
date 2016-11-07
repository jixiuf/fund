package eastmoney

import (
	"sort"
	"time"
)

var DefaultFetchTimeoutMS = 1000 * 20 // 20s

const (
	FundTypeAll   FundType = 1 // 所有开放式基金
	FundTypeStock FundType = 2 // 股票型开放式基金
	FundTypeBlend FundType = 3 // 混合型开放式基金
	FundTypeBond  FundType = 4 // 债券开放式基金
	FundTypeIndex FundType = 5 // 指数开放式基金
	FundTypeETF   FundType = 6 // ETF开放式基金
	FundTypeQDII  FundType = 7 // QDII开放式基金
	FundTypeLOF   FundType = 8 // LOF开放式基金
	FundTypeCNJY  FundType = 9 // 场内交易基金
)

type FundBase struct {
	Id   string
	Name string
	// ValueUnit           float64   // 单位净值
	// ValueTotal          float64   // 累计净值
	// LastValueUpdateTime time.Time // 净值最后一次更新日期
}

type Fund struct {
	FundBase
	Type string

	FundValueLast float64 // 最新一天的净值

	DayRatioLast            float64   // 最近一天的增长率
	FundValueLastUpdateTime time.Time // 净值的最后更新日期
	TotalFundValueLast      float64   // 最新累计净值

	// FundValueGuess               float64
	FundValueList        FundValueList // 净值列表
	TotalMoney           int64         // 基金规模
	TotalMoneyUpdateTime time.Time     // 基金规模更新时间
	MgrHeader            string        // 基金经理
	MgrHeaderId          string        // 基金经理id
	CreateTime           time.Time
}
type FenHongType int

const (
	FenHongType1 FenHongType = 1 // 1.每份基金份额折算1.012175663份
	FenHongType2 FenHongType = 2 // 2.每份派现金0.2150元,
	FenHongType3 FenHongType = 3 // 3. 每份基金份额分拆1.162668813份 (拆分后净值一般会变成1,用户持有份额会相应增加)
)

type FundValue struct {
	// 净值日期	单位净值	累计净值	日增长率	申购状态	赎回状态
	Value        float64 //
	TotalValue   float64
	DayRatio     float64
	Time         time.Time
	FenHongRatio float64     // 每份基金份额折算1.012175663份
	FenHongType  FenHongType // 1.每份基金份额折算1.012175663份 2.每份派现金0.2150元, 3. 每份基金份额分拆1.162668813份
}
type FundValueList []FundValue

func (l FundValueList) Sort() { // 按时间升序排列
	sort.Sort(l)
}

// 实现sort 接口
func (l FundValueList) Len() int {
	return len(l)
}
func (l FundValueList) Less(i, j int) bool {
	return l[i].Time.Before(l[j].Time)
}
func (l FundValueList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type FundList []Fund
type FundListSort struct {
	FundList     FundList
	sortFuncLess func(i, j Fund) bool
}

func NewFundListSort(fundList FundList) FundListSort {
	return FundListSort{FundList: fundList}

}

func (l *FundListSort) Sort(sortFuncLess func(i, j Fund) bool) { // 按时间升序排列
	l.sortFuncLess = sortFuncLess
	sort.Sort(l)
}

// 实现sort 接口
func (l *FundListSort) Len() int {
	return len(l.FundList)
}
func (l *FundListSort) Less(i, j int) bool {
	return l.sortFuncLess(l.FundList[i], l.FundList[j])
}
func (l *FundListSort) Swap(i, j int) {
	l.FundList[i], l.FundList[j] = l.FundList[j], l.FundList[i]
}
