package main

// 验证一个猜想，
// 假如最近一周净值一直是增长达3%，是否意味着接下来的一周也会继续增长

import (
	"fmt"
	"time"

	"github.com/jixiuf/fund/db"
	"github.com/jixiuf/fund/defs"
	"github.com/jixiuf/fund/dt"
	"github.com/jixiuf/fund/eastmoney"
)

const (
	YEAR time.Duration = time.Duration(time.Hour * 24 * 365)
	Day  time.Duration = time.Duration(time.Hour * 24)
)

var expectedRatio float64 = 0.03 // 期望的每周增幅

func main() {
	var result []float64
	dbT, _ := dt.NewDatabaseTemplateWithConfig(defs.DBConfig, true)

	fundList := db.FundValueHistoryGetAll(dbT)
	fListSorter := eastmoney.NewFundListSort(fundList)
	begin := time.Date(2017, 3, 6, 0, 0, 0, 0, time.Local)
	end := time.Date(2018, 3, 6, 0, 0, 0, 0, time.Local)
	from := begin
	for {
		to := from.Add(Day * 7)
		if to.After(end) {
			break
		}
		topNCandidates := filterTopN(&fListSorter, from, to, 300, expectedRatio)                  // 周增幅大于3.0的前300名做为候选fund
		nextWeekCandidates := filter(&fListSorter, from.Add(Day*7), to.Add(Day*7), expectedRatio) // 下一周增幅大于3.0
		cnt := 0
		for _, f := range topNCandidates {
			if nextWeekCandidates.IsIn(f) {
				cnt++
			}
		}
		if len(topNCandidates) != 0 {
			result = append(result, float64(cnt)/float64(len(topNCandidates)))
			fmt.Printf("%d.%d.%d-%d.%d.%d 下周达标基金=%d 候选基金总数=%d，下周达标基金/候选基金总数=%.2f%%\n",
				from.Year(), from.Month(), from.Day(),
				to.Year(), to.Month(), to.Day(),
				cnt, len(topNCandidates), 100*float64(cnt)/float64(len(topNCandidates)))
		}

		// 筛选出下一周 依然可以增符>3.0
		from = to

	}
}

// 获取收益率>minRatio 基金
func filter(l *eastmoney.FundListSort, from, to time.Time, minRatio float64) (ret eastmoney.FundList) {
	l.Sort(func(i, j eastmoney.Fund) bool {
		return i.CalcFundYield(from, to) > j.CalcFundYield(from, to)
	})
	for _, f := range l.FundList {
		if f.CalcFundYield(from, to) > minRatio {
			ret = append(ret, f)
		}

	}
	return

}

// 获取收益率>minRatio 的前topN名基金
func filterTopN(l *eastmoney.FundListSort, from, to time.Time, topN int, minRatio float64) (ret eastmoney.FundList) {
	l.Sort(func(i, j eastmoney.Fund) bool {
		return i.CalcFundYield(from, to) > j.CalcFundYield(from, to)
	})
	for i := 0; i < topN; i++ {
		if len(l.FundList) <= i {
			break
		}

		if l.FundList[i].CalcFundYield(from, to) > minRatio {
			// fmt.Println("", l.FundList[i].Id, l.FundList[i].Name, l.FundList[i].CalcFundYield(from, to))

			ret = append(ret, l.FundList[i])
		}
	}
	return

}
