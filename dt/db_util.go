package dt

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host         string `json:"host,omitempty"`
	User         string `json:"user,omitempty"`
	Pass         string `json:"passwd,omitempty"`
	Name         string `json:"database,omitempty"`
	Port         string `json:"port,omitempty"`
	MaxOpenConns int    `json:"max_active_connections,omitempty"`
	MaxIdleConns int    `json:"max_idle_connections,omitempty"`
	CharSet      string `json:"charset,omitempty"`
}
type MasterSlaveConfig struct {
	Master    DBConfig   `json:"master,omitempty"`
	SlaveList []DBConfig `json:"slave,omitempty"`
}

func (config MasterSlaveConfig) SlaveListLength() int {
	return len(config.SlaveList)
}

func ParseMasterSlaveConfig(jsonString string) (masterSlaveConfig MasterSlaveConfig, ok bool) {
	err := json.Unmarshal([]byte(jsonString), &masterSlaveConfig)
	if err != nil {
		fmt.Println("parse_master_slave_config_json_error", err)
		return
	}

	ok = true
	return

}

func NewDatabaseTemplateWithConfig(dbConfig DBConfig, debug bool) (dt DatabaseTemplate, ok bool) {
	var db *sql.DB
	db, ok = NewDBInstance(dbConfig, debug)
	if !ok {
		return
	}
	return &DatabaseTemplateImpl{db}, ok

}
func NewDatabaseTemplate(db *sql.DB) (dt DatabaseTemplate) {
	return &DatabaseTemplateImpl{db}
}
func NewDBInstance(dbConfig DBConfig, debug bool) (db *sql.DB, ok bool) {
	var (
		dbToken string
		err     error
		Log     string
	)

	if dbConfig.Port == "" {
		dbConfig.Port = "3306"
	}
	if dbConfig.CharSet == "" {
		dbConfig.CharSet = "utf8mb4"
	}

	dbToken = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local&tls=false&timeout=1m", dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.CharSet)
	Log = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local&tls=false&timeout=1m maxOpenConns=%d,maxIdleConns=%d\n", dbConfig.User, "password", dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.CharSet, dbConfig.MaxOpenConns, dbConfig.MaxIdleConns)
	db, err = sql.Open("mysql", dbToken)
	if err != nil {
		fmt.Println("error", Log, err)
		ok = false
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error", Log, err)
		ok = false
		return

	}
	if dbConfig.MaxOpenConns > 0 {
		db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}
	if dbConfig.MaxIdleConns > 0 {
		db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	}

	if debug {
		fmt.Print(Log)
	}
	ok = true
	return

}
