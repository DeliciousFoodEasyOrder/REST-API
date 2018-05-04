package models

import (
	_ "github.com/go-sql-driver/mysql" // Mysql Driver
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

const (
	dbDSN      string = "easyorder:Passw0r_@/easyorder"
	timeFormat string = "2006-01-02 15:04:05"
)

var (
	orm *xorm.Engine
)

func init() {
	var err error
	orm, err = xorm.NewEngine("mysql", dbDSN)
	orm.SetMapper(core.SameMapper{})
	if err != nil {
		panic(err)
	}
}
