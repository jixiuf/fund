package eastmoney

import (
	"bytes"
	"fmt"

	"bitbucket.org/jixiuf/fund/utils"
	"github.com/PuerkitoBio/goquery"
)

type FundType int

var DefaultFetchTimeoutMS = 1000 * 20 // 20s

const (
	FundTypeAll   FundType = 1 // 所有开放式基金
	FundTypeStock FundType = 2 // 股票型开放式基金
	FundTypeBlend FundType = 3 // 混合型开放式基金
	FundTypeBond  FundType = 4 // 债券开放式基金
	FundTypeIndex FundType = 5 // 指数开放式基金
	FundTypeETF   FundType = 6 // ETF开放式基金
	FundTypeQDII  FundType = 7 // QDII开放式基金
	FundTypeLOF   FundType = 8 // LOF开放式基金
	FundTypeCNJY  FundType = 9 // 场内交易基金
)

type FundBase struct {
	Id   string
	Name string
	// ValueUnit           float64   // 单位净值
	// ValueTotal          float64   // 累计净值
	// LastValueUpdateTime time.Time // 净值最后一次更新日期
}

const (
	EasyMoneyHome = "http://fund.eastmoney.com"
)

// 获得指定类型的开放式基金的 基金代码,一般是6位的数字形式如000011:华夏大盘精选
func GetFundIdList(openFuncType FundType) (list []FundBase) {
	urlStr := getUrlByFundType(openFuncType)
	if urlStr == "" {
		return
	}
	data, _ := utils.HttpGetWithReferer(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS)
	var doc *goquery.Document
	var err error
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("GetFundIdList.goquery.parse.error", err)
		return
	}

	doc.Find("#oTable tbody tr").Each(func(i int, tr *goquery.Selection) {
		var fs FundBase
		fs.Id = tr.Find("td").Eq(2).Text()
		fs.Name = tr.Find("td").Eq(3).Find("a").Eq(0).Text()
		list = append(list, fs)
	})
	return list

}
func getUrlByFundType(openFuncType FundType) string {
	switch openFuncType {
	case FundTypeAll:
		return "http://fund.eastmoney.com/fundguzhi.html"
	case FundTypeStock:
		return "http://fund.eastmoney.com/GP_fundguzhi3.html"
	case FundTypeBlend:
		return "http://fund.eastmoney.com/HH_fundguzhi4.html"
	case FundTypeBond:
		return "http://fund.eastmoney.com/ZQ_fundguzhi3.html"
	case FundTypeIndex:
		return "http://fund.eastmoney.com/ZS_fundguzhi3.html"
	case FundTypeETF:
		return "http://fund.eastmoney.com/ETF_fundguzhi3.html"
	case FundTypeQDII:
		return "http://fund.eastmoney.com/QDII_fundguzhi3.html"
	case FundTypeLOF:
		return "http://fund.eastmoney.com/LOF_fundguzhi3.html"
	case FundTypeCNJY:
		return "http://fund.eastmoney.com/cnjy_fundguzhi3.html"
	}

	return ""
}
