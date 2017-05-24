package main

import (
	"fmt"
	"time"

	"github.com/jixiuf/fund/db"
	"github.com/jixiuf/fund/defs"
	"github.com/jixiuf/fund/dt"
	"github.com/jixiuf/fund/eastmoney"
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
		db.FundValueHistoryUpdateType(dbT, fd)
		time.Sleep(time.Millisecond * 100)
	}

}
