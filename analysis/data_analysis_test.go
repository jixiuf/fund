package analysis

import "testing"

import "bitbucket.org/jixiuf/fund/eastmoney"
import "time"
import "github.com/stretchr/testify/assert"

func TestCalcFundYield(t *testing.T) {
	fd, err := eastmoney.GetFund("519033", true)
	assert.Nil(t, err)
	from := time.Date(2016, 8, 4, 0, 0, 0, 0, time.Local)
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)
	value := CalcFundYield(fd, from, to)
	assert.True(t, value-0.0531 < 0.0001, "海富通国策导向混合(519033) 2016-11-4的近3月收益率")
}
