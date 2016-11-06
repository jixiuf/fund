package main

import (
	"fmt"
	"time"

	"bitbucket.org/jixiuf/fund/db"
	"bitbucket.org/jixiuf/fund/defs"
	"bitbucket.org/jixiuf/fund/dt"
	"bitbucket.org/jixiuf/fund/eastmoney"
)

func main() {
	dailyUpdateeFundValue()
}
func dailyUpdateeFundValue() {
	dbT, _ := dt.NewDatabaseTemplateWithConfig(defs.DBConfig, true)
	db.FundValueHistoryCreateTable(dbT)
	stockList := eastmoney.GetFundIdList(eastmoney.FundTypeAll)
	for idx, fb := range stockList {
		fmt.Printf("%d/%d id=%s\n", idx, len(stockList), fb.Id)
		fd, err := eastmoney.GetFund(fb.Id, 10)
		if err != nil {
			fmt.Println(err)
			continue
		}
		db.FundValueHistoryInsertAll(dbT, fd) // 目前 无法处理 当日分红的情况, 待优化
		time.Sleep(time.Millisecond * 200)
	}

}
