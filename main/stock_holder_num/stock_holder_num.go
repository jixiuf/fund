package main

import (
	"fmt"
	"time"

	"github.com/jixiuf/fund/eastmoney"
	"github.com/tealeg/xlsx"
)

func main() {
	noticeDate, _ := time.ParseInLocation("2006-01-02", "2017-12-29", time.Local)
	list := eastmoney.GetStockHolderInfo()
	// list = list.RemoveByNoticeDate(noticeDate)
	list = list.RemoveByEndDate(noticeDate)
	list = list.FilterHolderNumChangeRate(-100, 0)
	list.Sort()
	if len(list) > 100 {
		list = list[:100]
	}
	fmt.Println(list)

	// csvFile, err := os.OpenFile("output.csv", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// w := csv.NewWriter(csvFile)
	f := xlsx.NewFile()
	sheet, _ := f.AddSheet("data")
	list.AddRows(sheet)
	f.Save("output.xlsx")

}
