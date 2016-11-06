package db

import (
	"database/sql"
	"fmt"

	"bitbucket.org/jixiuf/fund/dt"
	"bitbucket.org/jixiuf/fund/eastmoney"
)

func FundValueHistoryCreateTable(d dt.DatabaseTemplate) {
	sql := ` create table if not exists fund_value_history(
fundId varchar(6) default '',
name varchar(64) default '',
time timestamp default CURRENT_TIMESTAMP,
type varchar(64) default '' comment '类型',
value float default 0 comment '净值',
totalValue float default 0 comment '累计净值',
dayRatio float default 0 comment '日增比率',
fenHongType tinyint default 0 comment '1.每份基金份额折算1.012175663份 2.每份派现金0.2150元',
fenHongRatio float default 0 comment '分红比率,如每份基金份额折算1.012175663份',
primary key(fundId,time)
)`
	d.ExecDDL(sql)

}

func FundValueHistoryInsertAll(d dt.DatabaseTemplate, fd eastmoney.Fund) {
	sql := "insert into fund_value_history (fundId,name,type,time,value,totalValue,dayRatio,fenHongType,fenHongRatio) values"
	for idx, fv := range fd.FundValueList {
		sql += fmt.Sprintf("('%s','%s','%s','%s',%f,%f,%f,%d,%f)", fd.Id, fd.Name, fd.Type, fv.Time.Format("2006-01-02 00:00:00"), fv.Value, fv.TotalValue, fv.DayRatio, fv.FenHongType, fv.FenHongRatio)
		if idx != len(fd.FundValueList)-1 {
			sql += ",\n"
		}
	}
	sql += " ON DUPLICATE KEY UPDATE dayRatio=values(dayRatio),value=values(value),totalValue=values(totalValue),name=values(name),fenHongRatio=values(fenHongRatio),fenHongType=values(fenHongType)"
	d.ExecDDL(sql)
}

func FundValueHistoryInsertLast(d dt.DatabaseTemplate, fd eastmoney.Fund) {
	sql := "insert into fund_value_history (fundId,name,type,time,value,totalValue,dayRatio,fenHongRatio) values"
	sql += fmt.Sprintf("('%s','%s','%s','%s',%f,%f,%f,%f)", fd.Id, fd.Name, fd.Type, fd.FundValueLastUpdateTime.Format("2006-01-02 00:00:00"), fd.FundValueLast, fd.TotalFundValueLast, fd.DayRatioLast, 0.0)
	sql += " ON DUPLICATE KEY UPDATE dayRatio=values(dayRatio),type=values(type),value=values(value),totalValue=values(totalValue),name=values(name)"
	d.ExecDDL(sql)
}

func mapRowFundValueHistoryGetAll(resultSet *sql.Rows) (interface{}, error) {
	e := eastmoney.Fund{}
	history := eastmoney.FundValue{}

	err := resultSet.Scan(
		&e.Id,
		&e.Name,
		&e.Type,
		&history.Time,
		&history.Value,
		&history.TotalValue,
		&history.DayRatio,
		&history.FenHongRatio,
		&history.FenHongType,
	)
	e.FundValueList = append(e.FundValueList, history)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func FundValueHistoryGetAll(d dt.DatabaseTemplate) (fdList eastmoney.FundList) {
	sql := "select fundId,name,type,time,value,totalValue,dayRatio,fenHongRatio,fenHongType from fund_value_history order by fundId,time asc"
	e := eastmoney.Fund{}
	lastE := eastmoney.Fund{}
	arradd, err := d.QueryArray(sql, mapRowFundValueHistoryGetAll)
	if err != nil {
		return
	}
	for _, obj := range arradd {
		e = obj.(eastmoney.Fund)
		if e.Id == lastE.Id || lastE.Id == "" {
			lastE.FundValueList = append(lastE.FundValueList, e.FundValueList[0])
		} else {
			fdList = append(fdList, lastE)
			lastE = e
		}
	}

	return

}

func FundValueHistoryGet(d dt.DatabaseTemplate, fundId string) (fd eastmoney.Fund) {
	sql := fmt.Sprintf("select fundId,name,type,time,value,totalValue,dayRatio,fenHongRatio,fenHongType from fund_value_history where fundId=%sorder by fundId,time asc", fundId)
	e := eastmoney.Fund{}
	arradd, err := d.QueryArray(sql, mapRowFundValueHistoryGetAll)
	if err != nil {
		return
	}
	for _, obj := range arradd {
		e = obj.(eastmoney.Fund)
		fd.FundValueList = append(fd.FundValueList, e.FundValueList[0])
	}

	return

}
