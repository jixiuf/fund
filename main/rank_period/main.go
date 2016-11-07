package main

import (
	"fmt"

	"bitbucket.org/jixiuf/fund/db"
	"bitbucket.org/jixiuf/fund/defs"
	"bitbucket.org/jixiuf/fund/dt"
	"bitbucket.org/jixiuf/fund/eastmoney"
)

func main() {
	dbT, _ := dt.NewDatabaseTemplateWithConfig(defs.DBConfig, true)

	fundList := db.FundValueHistoryGetAll(dbT)
	fListSorter := eastmoney.NewFundListSort(fundList)
	fListSorter.Sort(func(i, j eastmoney.Fund) bool {
		a, _ := i.CalcFundPeriodYieldLastYear(eastmoney.Month)
		b, _ := j.CalcFundPeriodYieldLastYear(eastmoney.Month)
		return a > b
	})
	fmt.Println("近一年定投收益率排名")
	for idx, fd := range fListSorter.FundList {
		yield, count := fd.CalcFundPeriodYieldLastYear(eastmoney.Month)
		fmt.Println(idx, fd.Id, fd.Name, fd.Type, yield, count)
		if idx > 200 {
			break
		}
	}
}
