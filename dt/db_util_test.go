package dt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMasterSlaveConfig(t *testing.T) {
	jsonStr := `{"master":{"user":"th_dev","passwd":"th_devpass","database":"tapalliance_1_1","host":"localhost"},"slave":[]}`
	masterSlaveConfig, ok := ParseMasterSlaveConfig(jsonStr)
	assert.True(t, ok)
	assert.NotEmpty(t, masterSlaveConfig.Master.Host)
	fmt.Println(masterSlaveConfig.Master.Host)

}

func TestNewMasterSlaveConfig2(t *testing.T) {
	jsonStr := `{"master":{"user":"th_dev","passwd":"th_devpass","database":"tapalliance_1_1","host":"localhost"},"slave":[{"user":"th_dev","passwd":"th_devpass","database":"tapalliance_1_1","host":"localhost"},{"user":"th_dev","passwd":"th_devpass","database":"tapalliance_1_1","host":"localhost"}]}`
	masterSlaveConfig, ok := ParseMasterSlaveConfig(jsonStr)
	assert.True(t, ok)
	assert.NotEmpty(t, masterSlaveConfig.Master.Host)
	fmt.Println(masterSlaveConfig.Master.Host)
	assert.Equal(t, 2, len(masterSlaveConfig.SlaveList))

}
