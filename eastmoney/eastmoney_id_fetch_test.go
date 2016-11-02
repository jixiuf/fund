package eastmoney

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetFundIdList(t *testing.T) {
	list := GetFundIdList(FundTypeStock)
	assert.NotNil(t, list)
}
