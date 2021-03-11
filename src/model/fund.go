package model

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	. "oktools/src/global"
	"time"
)

type FundDataGenerated struct {
	Data       FundData `json:"data"`
	ResultCode int      `json:"result_code"`
}
type StockList struct {
	Name             string  `json:"name"`
	Code             string  `json:"code"`
	Percent          float64 `json:"percent"`
	CurrentPrice     float64 `json:"current_price"`
	ChangePercentage float64 `json:"change_percentage"`
	XqSymbol         string  `json:"xq_symbol"`
	XqURL            string  `json:"xq_url"`
	Amarket          bool    `json:"amarket"`
}
type FundPosition struct {
	StockPercent float64       `json:"stock_percent"`
	CashPercent  float64       `json:"cash_percent"`
	OtherPercent float64       `json:"other_percent"`
	AssetTot     float64       `json:"asset_tot"`
	AssetVal     float64       `json:"asset_val"`
	SourceMark   string        `json:"source_mark"`
	Source       string        `json:"source"`
	Enddate      string        `json:"enddate"`
	StockList    []StockList   `json:"stock_list"`
	BondList     []interface{} `json:"bond_list"`
}
type DeclareRateTable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type WithdrawRateTable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type OtherRateTable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type FundRates struct {
	FdCode            string              `json:"fd_code"`
	SubscribeRate     string              `json:"subscribe_rate"`
	DeclareRate       string              `json:"declare_rate"`
	WithdrawRate      string              `json:"withdraw_rate"`
	Discount          string              `json:"discount"`
	SubscribeDiscount string              `json:"subscribe_discount"`
	DeclareDiscount   string              `json:"declare_discount"`
	DeclareRateTable  []DeclareRateTable  `json:"declare_rate_table"`
	WithdrawRateTable []WithdrawRateTable `json:"withdraw_rate_table"`
	OtherRateTable    []OtherRateTable    `json:"other_rate_table"`
}
type AchievementList struct {
	FundCode  string  `json:"fund_code"`
	Fundsname string  `json:"fundsname"`
	PostDate  string  `json:"post_date"`
	CpRate    float64 `json:"cp_rate"`
	ResiDate  string  `json:"resi_date,omitempty"`
}
type ManagerList struct {
	Name            string            `json:"name"`
	Resume          string            `json:"resume"`
	College         string            `json:"college"`
	AchievementList []AchievementList `json:"achievement_list"`
}
type FundDateConf struct {
	FdCode          string `json:"fd_code"`
	BuyConfirmDate  int    `json:"buy_confirm_date"`
	BuyQueryDate    int    `json:"buy_query_date"`
	SaleConfirmDate int    `json:"sale_confirm_date"`
	SaleQueryDate   int    `json:"sale_query_date"`
	AllBuyDays      int    `json:"all_buy_days"`
	AllSaleDays     int    `json:"all_sale_days"`
}
type FundData struct {
	FundCompany  string        `json:"fund_company"`
	FundPosition FundPosition  `json:"fund_position"`
	FundRates    FundRates     `json:"fund_rates"`
	ManagerList  []ManagerList `json:"manager_list"`
	FundDateConf FundDateConf  `json:"fund_date_conf"`
}
type Fund struct {
	/*
	  `id` int(11) NOT NULL AUTO_INCREMENT,
	  `fund_name` varchar(64) DEFAULT '' COMMENT '基金名称',
	  `fund_code` varchar(16) NOT NULL DEFAULT '' COMMENT '基金代码',
	  `managers` varchar(32) NOT NULL DEFAULT '' COMMENT '管理人',
	  `enddate` varchar(32) NOT NULL DEFAULT '' COMMENT '季报日期',
	  `type` varchar(32) NOT NULL DEFAULT '' COMMENT '基金类型',
	  `detail_json` text NOT NULL COMMENT '蛋卷基金详细信息 json',
	  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
	  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
	  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
	*/

	FundName   string `json:"fund_name"`
	FundCode   string `json:"fund_code"`
	Managers   string `json:"managers"`
	Enddate    string `json:"enddate"`
	Type       string `json:"type"`
	DetailJson string `json:"detail_json"`

	ID        int64      `json:"id" gorm:"primary_key"`
	CreatedAt JSONTime   `json:"created_at" gorm:"column:created_at"` //创建日期
	UpdatedAt time.Time  `json:"-" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at" sql:"index"`

	DB *gorm.DB `json:"-" gorm:"-"` // 数据库操作DB

	FundDataShow FundDataGenerated `json:"fund_data_show" gorm:"-"`
}

/**
 * 数据打码
 */
func (obj *Fund) Init() {
	if obj == nil {
		return
	}
	err := json.Unmarshal([]byte(obj.DetailJson), &obj.FundDataShow)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

/**
 *  获取结构体，方便对结构体的数据进行初始化
 */
func GetFund() Fund {
	var obj Fund
	obj.DB = SqlDB.Conn
	fmt.Println(SqlDB.Database)
	return obj
}

// 设置User的表名为`profiles`
func (obj Fund) TableName() string {
	return "danjuan_fund"
}

// 获取数据库链接
func (obj Fund) getDB() (db *gorm.DB) {
	db = obj.DB
	if db == nil {
		db = SqlDB.GetConn()
	}
	db = db.Table(obj.TableName())
	db = db.Select(obj.FieldList())
	return db
}

func (obj Fund) FieldList() (fieldList []string) {
	fieldList = SqlDB.GetFieldName(obj)
	return
}

//传递过来的json数据初始化
func (obj *Fund) Insert() (id int64, err error) {
	err = obj.getDB().Create(&obj).Error
	if err == nil {
		id = obj.ID
	}
	return
}

// 更新部分, 可以清空
func (obj *Fund) Update(id int64, updateData map[string]interface{}) (err error) {
	err = obj.Get(id)
	if err != nil {
		return
	}
	updateData = SqlDB.CleanChangeMap(updateData, obj)
	err = obj.getDB().Model(&obj).Updates(updateData).Error
	return
}

//删除数据
func (obj *Fund) Delete(id int64) (err error) {
	err = obj.Get(id)
	if err != nil {
		return
	}
	err = obj.getDB().Delete(&obj).Error
	return
}

func (obj *Fund) GetCompanyByBankCardNum(company_id int64, bankNum string) (err error) {
	searchMap := make(map[string]string)
	searchMap["bd_num"] = bankNum

	err = obj.QueryOneByCondition(searchMap)
	return
}

//获取一个
func (obj *Fund) Get(id int64) (err error) {
	err = obj.getDB().First(&obj, id).Error
	if err == nil {
		obj.Init()
	}
	return
}

//获取查询db
func (obj *Fund) QueryDB(db *gorm.DB, searchKey map[string]string) (returnDB *gorm.DB) {

	// 查询扩展
	// 字段中如果字段 为 realname
	// 则默认支持以下查询
	// 		searchKey["realname"]  				精确查询		"realname = ?", "%" + searchKey["realname_min"] + "%"
	// 		searchKey["realname_not"]  			精确查询		"realname <> ?", "%" + searchKey["realname_min"] + "%"
	//		searchKey["realname_like"]  		模糊查询		"realname like ?", "%" + searchKey["realname_min"] + "%"
	//		searchKey["realname_left_like"]  	模糊左匹配查询		"realname like ?", searchKey["realname_min"] + "%"
	//		searchKey["realname_min"]			查询   		"realname >= ?" , searchKey["realname_min"]
	//		searchKey["realname_max"]			查询   		"realname <= ?" , searchKey["realname_min"]
	db = SqlDB.QueryDbInit(db, obj.FieldList(), searchKey)

	returnDB = db
	return
}

func (obj *Fund) QueryOneByCondition(searchKey map[string]string) (err error) {
	db := obj.getDB()
	db = obj.QueryDB(db, searchKey)
	err = db.First(&obj).Error

	if err == nil {
		obj.Init()
	}

	return
}

//列表
func (obj *Fund) List(searchKey map[string]string, page int64, pageSize int64) (results []Fund, err error) {
	db := obj.getDB()

	db = obj.QueryDB(db, searchKey)

	// 待调整成
	// select a.* from 表1 a, (select id from 表1 where 条件 limit 100000,20) b where a.id=b.id
	//

	// 分页
	if page > 0 && pageSize > 0 {
		if page <= 30 {
			//30页前直接获取
			db = db.Limit(pageSize).Offset((page - 1) * pageSize)
		} else {

			db = db.Select("id as b__id")
			db = db.Limit(pageSize).Offset((page - 1) * pageSize)

			db = obj.getDB().
				Table(obj.TableName()).
				Joins("join (?) b on id = b__id", db.QueryExpr())
		}
	} else {
		if pageSize > 0 {
			db = db.Limit(pageSize)
		}
	}

	// 查询数据
	err = db.Find(&results).Error

	if err == nil {
		for k, _ := range results {
			(&results[k]).Init()
		}
	}
	return
}

//获取数量
func (obj *Fund) Count(searchKey map[string]string) (totalSize int64, err error) {
	db := obj.getDB().Where("deleted_at is null")
	err = obj.QueryDB(db, searchKey).Count(&totalSize).Error

	return
}
