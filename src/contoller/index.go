package contoller

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"oktools/src/model"
	"time"
)

func FundIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "fund.html", gin.H{
		"times": GetTimes(),
		"funds": GetFunds(c),
	})
}

//
//func FundIndex(c *gin.Context) {
//	c.HTML(http.StatusOK, "index.html", gin.H{
//		//"tools": GetFunds(),
//	})
//}

func GetTimes() []struct{ Time string } {
	var times []struct{ Time string }
	now := time.Now()
	for i := 0; i < 120; i++ {
		date := now.AddDate(0, 0, -i)
		formatDay := date.Format("2006-01-02")
		times = append(times, struct{ Time string }{Time: formatDay})
	}
	return times
}

func GetFunds(c *gin.Context) []model.Fund {
	timeStr := c.Query("time")
	log.Print("GetFunds_Time", timeStr)
	searchMaps := make(map[string]string)
	searchMaps["enddate"] = timeStr
	searchMaps["order_by"] = "enddate desc"

	fund := model.GetFund()
	res, err := fund.List(searchMaps, 0, 100)
	if err != nil {
		return nil
	}

	return res
}
