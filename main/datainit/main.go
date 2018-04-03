package main

import (
	"fmt"
	"sync"

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
	// var tryAgainList []string
	var workerCnt = 100
	var queue [][]eastmoney.FundBase = make([][]eastmoney.FundBase, workerCnt)
	for idx, fb := range stockList {
		queue[idx%workerCnt] = append(queue[idx%workerCnt], fb)
	}
	var waitGroup sync.WaitGroup
	for _, q := range queue {
		waitGroup.Add(1)
		go func(q []eastmoney.FundBase) {
			for idx, fb := range q {
				fmt.Printf("%d/%d id=%s\n", idx, len(q), fb.Id)
				fd, err := eastmoney.GetFund(fb.Id, 0)
				if err != nil {
					fmt.Println(err)
					// tryAgainList = append(tryAgainList, fd.Id)
					continue
				}
				db.FundValueHistoryInsertAll(dbT, fd)
				time.Sleep(time.Millisecond * 10)
			}
			waitGroup.Done()

		}(q)

	}
	waitGroup.Wait()
	fmt.Println("failed_list", tryAgainList)

	// for idx, fundId := range tryAgainList {

	// 	fmt.Printf("%d/%d id=%s\n", idx, len(tryAgainList), fundId)
	// 	fd, err := eastmoney.GetFund(fundId, 0)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		tryAgainList = append(tryAgainList, fundId)
	// 		continue
	// 	}
	// 	db.FundValueHistoryInsertAll(dbT, fd)
	// 	time.Sleep(time.Millisecond * 10)
	// }

}
