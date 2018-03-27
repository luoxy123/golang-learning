package common

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DbContext struct {
	*gorm.DB
}

func NewDbContext(o MySqlOptions) *DbContext {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=true&loc=%s", o.UserName, o.Password, o.Host, o.Port, o.DataBase, o.Charset, url.QueryEscape(time.Local.String()))
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic("failed to connect database")
	}
	return &DbContext{db}
}
