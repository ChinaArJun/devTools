package contoller

import (
	"github.com/gin-gonic/gin"
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
	searchMaps := make(map[string]string)
	searchMaps[""] = c.Query("time")
	searchMaps[""] = ""
	searchMaps[""] = ""
	searchMaps[""] = ""

	fund := model.GetFund()
	res, err := fund.List(searchMaps, 0, 0)
	if err != nil {
		return nil
	}

	return res
}
