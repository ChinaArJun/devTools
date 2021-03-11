package contoller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
	"time"
)

func TestSql(t *testing.T) {
	db, err := gorm.Open("mysql", "root:amo1994@(127.0.0.1:3306)/python_fund?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := db.DB().Query("select * from danjuan_fund")

	var items []interface{}
	for rows.Next() {
		var countMap struct {
			Id int64 `json:"id"`
		}
		db.ScanRows(rows, &countMap)
		items = append(items, countMap)
	}
	rows.Close()
}

func TestGetTimes(t *testing.T) {
	now := time.Now()
	date := now.AddDate(0, 0, -1)
	formatDay := date.Format("2006-01-02")
	curMonthFormat := now.Format("2006-01")

	year, month, _ := now.Date()
	thisMonthFirstDay := time.Date(year, month, 1, 1, 1, 1, 0, now.Location())
	lastMonth := thisMonthFirstDay.AddDate(0, -1, 0)
	lastMonthFormat := lastMonth.Format("2006-01-02")

	lastday := thisMonthFirstDay.AddDate(0, 0, -1)
	fmt.Println(formatDay)
	fmt.Println(curMonthFormat)
	fmt.Println(lastMonthFormat)

	fmt.Println(lastday.Format("2006-01-02"))
}
