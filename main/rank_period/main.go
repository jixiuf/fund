package main

import (
	"fmt"

	"time"

	"github.com/jixiuf/fund/db"
	"github.com/jixiuf/fund/defs"
	"github.com/jixiuf/fund/dt"
	"github.com/jixiuf/fund/eastmoney"
)

// 计算定投排行,按

// | now~1yearAge | 1yearAge~2yearAge | 2yearAge~3YearAge | 3YearAge~5YearAge |
// |              |                   |                   |                   |
const (
	YEAR time.Duration = time.Duration(time.Hour * 24 * 365)
)

func main() {
	dbT, _ := dt.NewDatabaseTemplateWithConfig(defs.DBConfig, true)

	fundList := db.FundValueHistoryGetAll(dbT)
	fListSorter := eastmoney.NewFundListSort(fundList)
	period := eastmoney.Month
	var zeroTime time.Time
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	oneYearAge := now.Add(-YEAR)
	twoYearAge := now.Add(-(2 * YEAR))
	threeYearAge := now.Add(-(3 * YEAR))
	fiveYearAge := now.Add(-(5 * YEAR))

	thisYearYieldRank := getTopNFundPeriodYield(&fListSorter, period, oneYearAge, now, 100)          // 近一年定投收益前100名
	year2YieldRank := getTopNFundPeriodYield(&fListSorter, period, twoYearAge, oneYearAge, 100)      // 去年收益
	year3YieldRank := getTopNFundPeriodYield(&fListSorter, period, threeYearAge, twoYearAge, 100)    // 去去年收益
	year5_3YieldRank := getTopNFundPeriodYield(&fListSorter, period, fiveYearAge, threeYearAge, 100) // 3~5年之间 收益
	listSorter := thisYearYieldRank.And(year2YieldRank).And(year3YieldRank).And(year5_3YieldRank).GetSorter()
	listSorter.Sort(func(i, j eastmoney.Fund) bool { // 按照近5年收益率排名
		iallYield := i.CalcFundPeroidYield(period, fiveYearAge, now) //
		jallYield := j.CalcFundPeroidYield(period, fiveYearAge, now) //
		return iallYield > jallYield
	})
	fmt.Println("综合排名")
	for idx, fd := range listSorter.FundList {
		thisYearYield := fd.CalcFundPeroidYield(period, oneYearAge, now)
		Year2Yield := fd.CalcFundPeroidYield(period, twoYearAge, oneYearAge)
		Year3Yield := fd.CalcFundPeroidYield(period, threeYearAge, twoYearAge)
		Year5_3Yield := fd.CalcFundPeroidYield(period, fiveYearAge, threeYearAge)
		lastYear5Yield := fd.CalcFundPeroidYield(period, fiveYearAge, now)                      //
		allYield, allYieldCount := fd.CalcFundPeroidYieldWithPeriod(period, zeroTime, zeroTime) // 成立以来总收益,可投期数

		fmt.Printf("idx=%d,id=%s,name=%.20s,type=%.10s,this_year=%.4f,    1yearAge=%.4f,    2yearAge=%.4f,    3~5year=%.4f,last_5_year=%.4f all=%.04f,all_avg=%.4f\n",
			idx, fd.Id, fd.Name, fd.Type,
			thisYearYield, Year2Yield, Year3Yield, Year5_3Yield, lastYear5Yield, allYield, allYield/allYieldCount,
		)
	}
}
func print(list eastmoney.FundList) {

}

// 近一年定投收益率名前200
func getTopNFundPeriodYield(l *eastmoney.FundListSort, period eastmoney.Period, from, to time.Time, topN int) (ret eastmoney.FundList) {
	l.Sort(func(i, j eastmoney.Fund) bool {
		return i.CalcFundPeroidYield(period, from, to) > j.CalcFundPeroidYield(period, from, to)
	})
	if l.FundList.Len() > topN {
		ret = make(eastmoney.FundList, topN)
		copy(ret, l.FundList[:topN])
		return
	}
	ret = make(eastmoney.FundList, topN)
	copy(ret, l.FundList)
	return

	// for _, fd := range l.FundList {
	// 	_, count := fd.CalcFundPeriodYieldLastYear(period)
	// 	if count == 12 {
	// 		ret = append(ret, fd)
	// 	}
	// 	if len(ret) == topN {
	// 		return
	// 	}
	// }

	// return
}
