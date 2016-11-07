package eastmoney

import "testing"

import "time"
import "github.com/stretchr/testify/assert"

func TestCalcFundYield(t *testing.T) {
	fd, err := GetFund("519033", 0)
	assert.Nil(t, err)
	from := time.Date(2016, 8, 4, 0, 0, 0, 0, time.Local)
	to := time.Date(2016, 11, 4, 0, 0, 0, 0, time.Local)
	value := fd.CalcFundYield(from, to)
	assert.True(t, value-0.0531 < 0.0001, "海富通国策导向混合(519033) 2016-11-4的近3月收益率")
}
