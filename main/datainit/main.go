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
	initFundValueHistory()
}

// 初始化历史净值数据
func initFundValueHistory() {
	dbT, _ := dt.NewDatabaseTemplateWithConfig(defs.DBConfig, true)
	db.FundValueHistoryCreateTable(dbT)
	stockList := eastmoney.GetFundIdList(eastmoney.FundTypeAll)
	var tryAgainList []string
	for idx, fb := range stockList {
		fmt.Printf("%d/%d id=%s\n", idx, len(stockList), fb.Id)
		fd, err := eastmoney.GetFund(fb.Id, 0)
		if err != nil {
			fmt.Println(err)
			tryAgainList = append(tryAgainList, fd.Id)
			continue
		}
		db.FundValueHistoryInsertAll(dbT, fd)
		time.Sleep(time.Second)
	}
	fmt.Println("failed_list", tryAgainList)

	for idx, fundId := range tryAgainList {

		fmt.Printf("%d/%d id=%s\n", idx, len(tryAgainList), fundId)
		fd, err := eastmoney.GetFund(fundId, 0)
		if err != nil {
			fmt.Println(err)
			tryAgainList = append(tryAgainList, fundId)
			continue
		}
		db.FundValueHistoryInsertAll(dbT, fd)
		time.Sleep(time.Second)
	}

}
