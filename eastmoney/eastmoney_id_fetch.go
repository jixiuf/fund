package eastmoney

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

import "bitbucket.org/jixiuf/fund/utils"

type OpenFundType int

var DefaultFetchTimeoutMS = 1000 * 20 // 20s

const (
	OpenFundTypeAll   OpenFundType = 1 // 所有开放式基金
	OpenFundTypeStock OpenFundType = 2 // 股票型开放式基金
	OpenFundTypeBlend OpenFundType = 3 // 混合型开放式基金
	OpenFundTypeBond  OpenFundType = 4 // 债券开放式基金
	OpenFundTypeIndex OpenFundType = 5 // 指数开放式基金
	OpenFundTypeETF   OpenFundType = 6 // ETF开放式基金
	OpenFundTypeQDII  OpenFundType = 7 // QDII开放式基金
	OpenFundTypeLOF   OpenFundType = 8 // LOF开放式基金
)

type OpenFuncSimple struct {
	Id                  string
	Name                string
	ValueUnit           float64   // 单位净值
	ValueTotal          float64   // 累计净值
	LastValueUpdateTime time.Time // 净值最后一次更新日期
}

const (
	EasyMoneyHome = "http://fund.eastmoney.com"
)

// 获得指定类型的开放式基金的 基金代码,一般是6位的数字形式如000011:华夏大盘精选
func GetOpenFundIdList(openFuncType OpenFundType) (list []OpenFuncSimple) {
	urlStr := getUrlByOpenFundType(openFuncType)
	if urlStr == "" {
		return
	}
	fmt.Println(urlStr)
	data, _ := utils.HttpGetWithReferer(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS)
	var doc *goquery.Document
	var err error
	// fmt.Println(string(data))
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("GetOpenFundIdList.goquery.parse.error", err)
		return
	}

	tBody := doc.Find("#oTable tbody tr").Each(func(i int, tr *goquery.Selection) {
		fmt.Println(i, tr.Text())

	})
	fmt.Println(tBody)
	return nil

}
func getUrlByOpenFundType(openFuncType OpenFundType) string {
	switch openFuncType {
	case OpenFundTypeAll:
		return "http://fund.eastmoney.com/jzzzl.html#os_0;isall_1;ft_;pt_1"
	case OpenFundTypeStock:
		return "http://fund.eastmoney.com/GP_jzzzl.html#os_0;isall_1;ft_;pt_2"
	case OpenFundTypeBlend:
	case OpenFundTypeBond:
	case OpenFundTypeIndex:
	case OpenFundTypeETF:
	case OpenFundTypeQDII:
	case OpenFundTypeLOF:

	}

	return ""
}
