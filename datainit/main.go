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
	db.FundValueHistoryCreateTable(dbT)
	stockList := eastmoney.GetFundIdList(eastmoney.FundTypeStock)
	for _, fb := range stockList {
		fd, err := eastmoney.GetFund(fb.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		db.FundValueHistoryInsertAll(dbT, fd)
		break
	}
}
