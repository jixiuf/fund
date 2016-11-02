package eastmoney

import (
	"fmt"
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"

	"bytes"

	"strings"

	"bitbucket.org/jixiuf/fund/utils"
)

type Fund struct {
	FundBase
	FundValueList FundValueList
	TotalMoney    int64
}
type FundValue struct {
	// 净值日期	单位净值	累计净值	日增长率	申购状态	赎回状态
	Value      float64 //
	TotalValue float64
	DayRatio   float64
	Time       time.Time
}
type FundValueList []FundValue

func (l FundValueList) Sort() {
	sort.Sort(l)
}

// 实现sort 接口
func (l FundValueList) Len() int {
	return len(l)
}
func (l FundValueList) Less(i, j int) bool {
	return l[i].Time.Before(l[j].Time)
}
func (l FundValueList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func GetFund(fundId string) (f Fund, err error) {
	f, err = GetFundDetail(fundId)
	if err != nil {
		return
	}
	list, err := GetFundHistoryValueList(fundId)
	if err != nil {
		return f, err
	}
	f.FundValueList = list
	return
}

func GetFundDetail(fundId string) (f Fund, err error) {
	urlStr := fmt.Sprintf("http://fund.eastmoney.com/%s.html", fundId)
	data, err := utils.HttpGetWithReferer(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS)
	if err != nil {
		return f, err
	}

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("GetFundDetail.goquery.parse.error", err)
		return
	}

	// 信诚中证800有色指数分级(165520)
	title := doc.Find("div.fundDetail-tit div").Eq(0).Text()
	lastIndex := strings.LastIndex(title, "(")
	if lastIndex != -1 {
		title = title[:lastIndex]
	}

	f.FundBase.Name = title
	return
}

// http://fund.eastmoney.com/f10/jjjz_165520.html

//返回的内容是如下的一段javascript脚本，
// 要做的就是做content内容挖出来，
// var apidata={ content:"<table class='w782 comm lsjz'><thead><tr><th class='first'>净值日期</th><th>单位净值</th><th>累计净值</th><th>日增长率</th><th>申购状态</th><th>赎回状态</th><th class='tor last'>分红送配</th></tr></thead>
// <tbody><tr><td>2016-11-01</td><td class='tor bold'>1.0330</td><td class='tor bold'>1.1830</td><td class='tor bold red'>0.78%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr>
// </tbody></table>",records:764,pages:1,curpage:1};

func GetFundHistoryValueList(fundId string) (list FundValueList, err error) {
	// type=lsjz 历史净值
	perPage := 10000000
	urlStr := fmt.Sprintf("http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=%s&page=1&per=%d", fundId, perPage)
	data, err := utils.HttpGetWithReferer(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS)
	if err != nil {
		return nil, err
	}
	start := []byte(`content:"`)
	end := []byte(`</table>`)
	startPos := bytes.Index(data, start)
	endPos := bytes.Index(data, end)
	fmt.Println(startPos, endPos)
	if startPos == -1 || endPos == -1 {
		return nil, err
	}
	// <table class='w782 comm lsjz'><thead><tr><th class='first'>净值日期</th><th>单位净值</th><th>累计净值</th><th>日增长率</th><th>申购状态</th><th>赎回状态</th><th class='tor last'>分红送配</th></tr></thead><tbody><tr><td>2016-11-01</td><td class='tor bold'>1.0330</td><td class='tor bold'>1.1830</td><td class='tor bold red'>0.78%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr></tbody></table>
	content := data[startPos+len(start) : endPos+len(end)]

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		fmt.Println("GetFundHistoryValueList.goquery.parse.error", err)
		return
	}

	doc.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		var fv FundValue
		fv.Value = utils.Str2Float64(tr.Find("td").Eq(1).Text(), 0)
		fv.TotalValue = utils.Str2Float64(tr.Find("td").Eq(2).Text(), 0)
		fv.DayRatio = utils.Str2Float64(tr.Find("td").Eq(3).Text(), 0)
		fv.Time, _ = time.ParseInLocation("2006-01-02", tr.Find("td").Eq(0).Text(), time.Local)
		list = append(list, fv)
	})
	list.Sort()

	return
}
