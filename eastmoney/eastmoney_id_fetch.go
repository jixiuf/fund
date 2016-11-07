package eastmoney

import (
	"bytes"
	"fmt"

	"bitbucket.org/jixiuf/fund/utils"
	"github.com/PuerkitoBio/goquery"
)

type FundType int

func (openFuncType FundType) String() string {
	switch openFuncType {
	case FundTypeAll:
		return "All"
	case FundTypeStock:
		return "股票型"
	case FundTypeBlend:
		return "混合"
	case FundTypeBond:
		return "债券"
	case FundTypeIndex:
		return "股票指数"
	case FundTypeETF:
		return "ETF"
	case FundTypeQDII:
		return "QDII"
	case FundTypeLOF:
		return "LOF"
	case FundTypeCNJY:
		return "场内交易"
	}
	return ""

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
	data, _ := utils.HttpGetWithRefererTryN(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS, 5)
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
