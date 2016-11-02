package eastmoney

import (
	"fmt"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestGetFundHistoryValueList(t *testing.T) {
	list, err := GetFundHistoryValueList("165520")
	assert.Nil(t, err)
	assert.NotNil(t, list)
}

func TestGetFundDetail(t *testing.T) {
	list, err := GetFundDetail("165520")
	assert.Nil(t, err)
	fmt.Println(list)
}
