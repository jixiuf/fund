package db

import (
	"fmt"

	"bitbucket.org/jixiuf/fund/dt"
	"bitbucket.org/jixiuf/fund/eastmoney"
)

func FundValueHistoryCreateTable(d dt.DatabaseTemplate) {
	sql := ` create table if not exists fund_value_history(
fundId varchar(6) default '',
time timestamp default CURRENT_TIMESTAMP,
value float default 0,
totalValue float default 0,
dayRatio float default 0,
primary key(fundId,time)
)`
	d.ExecDDL(sql)

}

func FundValueHistoryInsertAll(d dt.DatabaseTemplate, fd eastmoney.Fund) {
	sql := "insert into fund_value_history (fundId,time,value,totalValue,dayRatio) values("
	for idx, fv := range fd.FundValueList {
		sql += fmt.Sprintf("'%s','%s',%f,%f,%f", fd.Id, fv.Time.Format("2006-01-02 00:00:00"), fv.Value, fv.TotalValue, fv.DayRatio)
		if idx != len(fd.FundValueList)-1 {
			sql += ",\n"
		}
	}
	sql += ")"
	fmt.Println(sql)

}
