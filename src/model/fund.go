package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	. "oktools/src/global"
	"strconv"
	"time"
)

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

	CpcnpaySrcsrl       string `json:"cpcnpay_srcsrl"`        // 中金支付绑卡交易流水号
	CpcnpayTxsn         string `json:"cpcnpay_txsn"`          // 中金绑卡平台方交易流水号
	CpcnpayVerifyStatus int64  `json:"cpcnpay_verify_status"` // 中金绑卡验证状态（0待验证 1通过 2失效）
	CpcnpayDefault      int64  `json:"cpcnpay_default"`       // 中金是否默认银行卡（0否 1是）

	ID        int64      `json:"id" gorm:"primary_key"`
	CreatedAt JSONTime   `json:"created_at" gorm:"column:created_at"` //创建日期
	UpdatedAt time.Time  `json:"-" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at" sql:"index"`

	//方便做复杂数据及复杂逻辑 可以做数据回滚
	DB *gorm.DB `json:"-" gorm:"-"` // 数据库操作DB

	VerifyStatusShow string `json:"verify_status_show" gorm:"-"` //验证状态中文
}

/**
 * 数据打码
 */
func (obj *Fund) Init() {
	if obj == nil {
		return
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

//根据公司id查询账户
func (obj *Fund) GetFundByCompanyId(companyId int64) (err error) {
	err = obj.getDB().Where("company_id = ?", companyId).First(&obj).Error

	if err == nil {
		obj.Init()
	}
	return
}

func (obj *Fund) GetByBankCardNum(bankNum string) (err error) {
	err = obj.getDB().Where("bank_card_num = ?", bankNum).First(&obj).Error
	if err == nil {
		obj.Init()
	}
	return
}

func (obj *Fund) GetCompanyByBankCardNum(company_id int64, bankNum string) (err error) {
	searchMap := make(map[string]string)
	searchMap["bank_card_num"] = bankNum

	err = obj.QueryOneByCondition(searchMap)
	return
}

func (obj *Fund) GetFundByDepositBank(company_id, verify_status int64, deposit_bank string) (err error) {
	searchMap := make(map[string]string)
	searchMap["company_id"] = strconv.FormatInt(company_id, 10)
	searchMap["verify_status"] = strconv.FormatInt(verify_status, 10)
	searchMap["deposit_bank_like"] = deposit_bank

	err = obj.QueryOneByCondition(searchMap)

	return
}

func (obj *Fund) GetByCfcaBankCardNum(bankNum string, verifyStatus int64) (err error) {
	searchMap := make(map[string]string)
	searchMap["bank_card_num"] = bankNum
	searchMap["cpcnpay_verify_status"] = strconv.FormatInt(verifyStatus, 10)
	searchMap["cpcnpay_srcsrl_is_not_null"] = "notNull"
	searchMap["cpcnpay_srcsrl_is_not_empty"] = "notEmpty"

	err = obj.QueryOneByCondition(searchMap)

	return
}

func (obj *Fund) GetCfcacount(companyId, verifyStatus int64) (err error) {
	searchMap := make(map[string]string)
	searchMap["company_id"] = strconv.FormatInt(companyId, 10)
	searchMap["cpcnpay_verify_status"] = strconv.FormatInt(verifyStatus, 10)
	searchMap["cpcnpay_srcsrl_is_not_null"] = "notNull"
	searchMap["cpcnpay_srcsrl_is_not_empty"] = "notEmpty"

	err = obj.QueryOneByCondition(searchMap)

	return
}

func (obj *Fund) GetCfcaDefaultAcount(companyId, verifyStatus, defaultStatus int64) (err error) {
	searchMap := make(map[string]string)
	searchMap["company_id"] = strconv.FormatInt(companyId, 10)
	searchMap["cpcnpay_verify_status"] = strconv.FormatInt(verifyStatus, 10)
	searchMap["cpcnpay_default"] = strconv.FormatInt(defaultStatus, 10)
	searchMap["cpcnpay_srcsrl_is_not_null"] = "notNull"
	searchMap["cpcnpay_srcsrl_is_not_empty"] = "notEmpty"

	err = obj.QueryOneByCondition(searchMap)

	return
}

func (obj *Fund) GetFundBankList(company_id, verify_status int64) (result []Fund, err error) {
	searchMap := make(map[string]string)
	searchMap["company_id"] = strconv.FormatInt(company_id, 10)
	searchMap["cpcnpay_verify_status"] = strconv.FormatInt(verify_status, 10)
	searchMap["order_by"] = "id asc"

	result, err = obj.List(searchMap, 0, 0)

	return
}

func (obj *Fund) GetCompanyFastPayAccountBankList(company_id, verify_status int64) (result []Fund, err error) {
	searchMap := make(map[string]string)
	searchMap["company_id"] = strconv.FormatInt(company_id, 10)
	searchMap["cpcnfastpay_verify_status"] = strconv.FormatInt(verify_status, 10)
	searchMap["order_by"] = "id asc"

	result, err = obj.List(searchMap, 0, 0)

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
