package eastmoney

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"

	"bytes"

	"strings"

	"github.com/jixiuf/fund/utils"
)

//
func GetFund(fundId string, fetchFundValueHistoryCnt int) (f Fund, err error) {
	f, err = GetFundDetail(fundId)
	if err != nil {
		return
	}
	if fetchFundValueHistoryCnt != -1 {
		list, err := GetFundHistoryValueList(fundId, fetchFundValueHistoryCnt)
		if err != nil {
			return f, err
		}
		f.FundValueList = list

	}
	return
}

func GetFundDetail(fundId string) (f Fund, err error) {
	f.Id = fundId
	urlStr := fmt.Sprintf("http://fund.eastmoney.com/%s.html", fundId)
	data, err := utils.HttpGetWithRefererTryN(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS, 5)
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

	// 基金类型
	fundType := doc.Find("div.infoOfFund table tbody tr").Eq(0).Find("td a").Eq(0).Text()
	f.Type = fundType

	// 基金规模
	// 基金规模：2.24亿元（2016-09-30）
	totalMoney := doc.Find("div.infoOfFund table tbody tr").Eq(0).Find("td").Eq(1).Text()
	var ft float64
	var totalMoneyUpdateTime string
	fmt.Sscanf(totalMoney, "基金规模：%f亿元（%10s）", &ft, &totalMoneyUpdateTime)
	f.TotalMoney = int64(ft * 10000 * 10000)
	f.TotalMoneyUpdateTime, _ = time.ParseInLocation("2006-01-02", totalMoneyUpdateTime, time.Local)

	//
	mgrHeaderNode := doc.Find("div.infoOfFund table tbody tr").Eq(0).Find("td a").Eq(2)
	f.MgrHeader = mgrHeaderNode.Text()
	// http://fund.eastmoney.com/f10/jjjl_002207.html
	href, ok := mgrHeaderNode.Attr("href")
	if ok {
		startIdx := strings.Index(href, "jjjl_")
		endIdx := strings.Index(href, ".html")
		f.MgrHeaderId = href[startIdx+len("jjjl_") : endIdx]
	}

	//
	createDate := doc.Find("div.infoOfFund table tbody tr").Eq(1).Find("td").Eq(0).Text()
	var createDateStr string
	fmt.Sscanf(createDate, "成 立 日：%s", &createDateStr)
	f.CreateTime, _ = time.ParseInLocation("2006-01-02", createDateStr, time.Local)

	// 最新净值
	fundValueLastNode := doc.Find("div.dataOfFund dl.dataItem02")
	// 单位净值 (2016-11-02)
	var fundValueLastUpdateTimeStr string
	fmt.Sscanf(fundValueLastNode.Find("dt").Text(), "单位净值 (%10s)", &fundValueLastUpdateTimeStr)
	f.FundValueLastUpdateTime, _ = time.ParseInLocation("2006-01-02", fundValueLastUpdateTimeStr, time.Local)
	f.FundValueLast = utils.Str2Float64(fundValueLastNode.Find("dd.dataNums span").Eq(0).Text(), 0)
	dayRatioLastStr := fundValueLastNode.Find("dd.dataNums span").Eq(1).Text()
	if dayRatioLastStr != "" {
		f.DayRatioLast = utils.Str2Float64(dayRatioLastStr[:len(dayRatioLastStr)-1], 0)
	}

	f.TotalFundValueLast = utils.Str2Float64(doc.Find("div.dataOfFund dl.dataItem03 dd.dataNums span").Text(), 0)

	return
}

// http://fund.eastmoney.com/f10/jjjz_165520.html

//返回的内容是如下的一段javascript脚本，
// 要做的就是做content内容挖出来，
// var apidata={ content:"<table class='w782 comm lsjz'><thead><tr><th class='first'>净值日期</th><th>单位净值</th><th>累计净值</th><th>日增长率</th><th>申购状态</th><th>赎回状态</th><th class='tor last'>分红送配</th></tr></thead>
// <tbody><tr><td>2016-11-01</td><td class='tor bold'>1.0330</td><td class='tor bold'>1.1830</td><td class='tor bold red'>0.78%</td><td>开放申购</td><td>开放赎回</td><td class='red unbold'></td></tr>
// </tbody></table>",records:764,pages:1,curpage:1};

func GetFundHistoryValueList(fundId string, perPage int) (list FundValueList, err error) {
	// type=lsjz 历史净值
	if perPage == 0 {
		perPage = 1000000000
	}

	urlStr := fmt.Sprintf("http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=%s&page=1&per=%d", fundId, perPage)
	data, err := utils.HttpGetWithRefererTryN(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS, 5)
	if err != nil {
		return nil, err
	}
	start := []byte(`content:"`)
	end := []byte(`</table>`)
	startPos := bytes.Index(data, start)
	endPos := bytes.Index(data, end)
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
		dayRatio := tr.Find("td").Eq(3).Text() // 0.15%
		if dayRatio != "" {
			fv.DayRatio = utils.Str2Float64(dayRatio[0:len(dayRatio)-1], 0)
		}

		fv.Time, _ = time.ParseInLocation("2006-01-02", tr.Find("td").Eq(0).Text(), time.Local)

		// 每份基金份额折算1.012175663份
		// 是否分红了
		fenHong := tr.Find("td").Eq(6).Text()
		if fenHong != "" {
			_, err := fmt.Sscanf(fenHong, "每份基金份额折算%f份", &fv.FenHongRatio)
			fv.FenHongType = FenHongType1
			if err != nil {
				_, err := fmt.Sscanf(fenHong, "每份派现金%f元", &fv.FenHongRatio)
				fv.FenHongType = FenHongType2
				if err != nil {
					fv.FenHongType = FenHongType3
					_, err := fmt.Sscanf(fenHong, "每份基金份额分拆%f份", &fv.FenHongRatio)
					if err != nil {
						fmt.Println("解析基金分红error:", fenHong, fundId, err)
					}
				}
			}
		}

		list = append(list, fv)
	})
	list.Sort()

	return
}
