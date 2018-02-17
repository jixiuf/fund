package eastmoney

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/jixiuf/fund/utils"
)

// 每个季度大都会公布财报信息，获取所有股票的股东数 变化信息，
// 想实际的功能， 如果股东数减少，说明筹码分布趋向于集中，
// 在公布财报时 将所有股票按股东减少比例排序，然后以此作为选股最基本的依据
// 但是天天基金网上，排序时，有些股票当季数据还没更新，也会排序靠前，所以参考意义不大
// 本工具会给定一个时间点，只有晚于此时间点的变动才会考虑
// http://data.eastmoney.com/gdhs/
// http://data.eastmoney.com/DataCenter_V3/gdhs/GetList.ashx?pagesize=4000&page=1
// "SecurityCode":"600165","SecurityName":"新日恒力","LatestPrice":"0","PriceChangeRate":"0","HolderNum":"15338","PreviousHolderNum":"19153",
// "HolderNumChange":"-3815","HolderNumChangeRate":"-19.9186","RangeChangeRate":"13.39","EndDate":"2018-02-12T00:00:00",
// "PreviousEndDate":"2017-09-30T00:00:00","HolderAvgCapitalisation":"737216.790014995",
// "HolderAvgStockQuantity":"44652.74","TotalCapitalisation":"11307431125.25","CapitalStock":"684883775","NoticeDate":"2018-02-15T00:00:00"}
type StockHolderInfo struct {
	Code                string       `json:"SecurityCode,omitempty"`
	Name                string       `json:"SecurityName,omitempty"`
	HolderNum           int          `json:"HolderNum,string,omitempty"`         // 本次财报公布股东数
	PreviousHolderNum   int          `json:"PreviousHolderNum,string,omitempty"` // 上期财报公布股东数
	HolderNumChangeRate json.Number  `json:"HolderNumChangeRate,omitempty"`      // 增减比例(%)
	RangeChangeRate     json.Number  `json:"RangeChangeRate,omitempty"`          // 这段期间 股价涨跌比例
	NoticeDate          JsonDateTime `json:"NoticeDate,omitempty"`               //公布日期
	PreviousEndDate     JsonDateTime `json:"PreviousEndDate,omitempty"`          //上期 统计日期
	EndDate             JsonDateTime `json:"EndDate,omitempty"`                  //本期 统计日期
}

func (this StockHolderInfo) String() string {
	return fmt.Sprintf("code: %s,name:%s,股东增少比例:%s,本期日期:%s,上期日期:%s,公布日期:%s,本期股东数:%d,上期股东数:%d,区间股价波动:%s\n",
		this.Code, this.Name, this.HolderNumChangeRate,
		this.EndDate.String(),
		this.PreviousEndDate.String(),
		this.NoticeDate.String(),
		this.HolderNum, this.PreviousHolderNum,
		this.RangeChangeRate)
}

type StockHolderInfoList []StockHolderInfo

func (l StockHolderInfoList) RemoveByNoticeDate(date time.Time) (l2 StockHolderInfoList) {
	l2 = make(StockHolderInfoList, 0, len(l))
	// 删除公布日期小于指定日期的数据
	for _, info := range l {
		if info.NoticeDate.ToTime().Before(date) {
			continue
		}
		l2 = append(l2, info)
	}
	return
}

func (l StockHolderInfoList) Sort() {
	sort.Sort(l)
}

// 实现sort 接口
func (l StockHolderInfoList) Len() int {
	return len(l)
}

func (l StockHolderInfoList) Less(i, j int) bool {
	iv, _ := l[i].HolderNumChangeRate.Float64()
	jv, _ := l[j].HolderNumChangeRate.Float64()
	return iv < jv
}

func (l StockHolderInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type stockHolkerInfoResponse struct {
	Success bool                `json:"success,omitempty"`
	Data    StockHolderInfoList `json:"data,omitempty"`
}

func GetStockHolderInfo() (list StockHolderInfoList) {
	// 4000>大于股票的总个数即可，这里写死了， 一页数据全拉下来
	urlStr := fmt.Sprintf("http://data.eastmoney.com/DataCenter_V3/gdhs/GetList.ashx?pagesize=4000&page=1")
	data, err := utils.HttpGetWithRefererTryN(urlStr, EasyMoneyHome, DefaultFetchTimeoutMS, 5)
	if err != nil {
		fmt.Println("GetStockHolderInfo", err)
		return
	}
	var response stockHolkerInfoResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		fmt.Println(err)
	}

	return response.Data

}
