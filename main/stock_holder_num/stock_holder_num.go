package main

import (
	"fmt"
	"time"

	"github.com/jixiuf/fund/eastmoney"
)

func main() {
	noticeDate, _ := time.ParseInLocation("2006-01-02", "2018-01-01", time.Local)
	list := eastmoney.GetStockHolderInfo()
	list = list.RemoveByNoticeDate(noticeDate)
	list.Sort()
	fmt.Println(list)

}
