package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	. "oktools/src/global"
	"time"
)

// 内存缓存有效期
const toolsCacheAge = 24 * time.Hour

type Tool struct {
	Path       string   `json:"path"`
	Title      string   `json:"title"`
	Icon       *string  `json:"icon"`
	Category   int      `json:"category,,omitempty"`
	UsageCount int      `json:"usage_count,omitempty"`
	DB         *gorm.DB `json:"-" gorm:"-"` // 数据库操作DB
}

type Cache struct {
	Data  interface{}
	Expir time.Time
}

var ToolMap = make(map[string]Tool)

// 全部工具内存缓存
var toolsCache *Cache

func GetTools() []Tool {
	now := time.Now()
	if toolsCache == nil || toolsCache.Expir.Before(now) {
		tools, err := loadTools()
		if err != nil {
			log.Println("Failed to load tools data :", err)
			return nil
		}

		toolsCache = &Cache{}
		toolsCache.Data = tools
		toolsCache.Expir = now.Add(toolsCacheAge)
		return tools
	}
	return toolsCache.Data.([]Tool)
}

func loadTools() ([]Tool, error) {
	//rows, err := global.PSDB.Query(`SELECT path,title,icon,usage_count,category FROM tools ORDER  BY usage_count desc`)
	//if err != nil {
	//	return nil, err
	//}
	//rows, err :=

	ts := GetTool()
	tools, _ := ts.List(nil, 0, 0)
	//for rows.Next() {
	//	tool := &Tool{}
	//	err = rows.Scan(&tool.Path, &tool.Title, &tool.Icon, &tool.UsageCount, &tool.Category)
	//	if err != nil {
	//		return nil, err
	//	}
	//	tools = append(tools, tool)
	//	ToolMap[tool.Path] = tool
	//}
	return tools, nil
}

func UpdateUsageCount() {
	//if len(ToolMap) == 0 {
	//	log.Println("Tools map is empty , No need to update data!")
	//	return
	//}
	//
	//var buf bytes.Buffer
	//buf.WriteString("UPDATE tools SET usage_count = CASE path ")
	//for k, v := range ToolMap {
	//	buf.WriteString(" WHEN '")
	//	buf.WriteString(k)
	//	buf.WriteString("' THEN ")
	//	buf.WriteString(strconv.Itoa(v.UsageCount))
	//}
	//buf.WriteString(" END WHERE path IN(")
	//
	//for k := range ToolMap {
	//	buf.WriteString("'")
	//	buf.WriteString(k)
	//	buf.WriteString("'")
	//	buf.WriteString(",")
	//}
	//buf.WriteString("'')")
	//
	//sql := buf.String()
	//_, err := PSDB.Exec(sql)
	//if err != nil {
	//	log.Println("Failed to update tools usage count data :", err)
	//}
	//
	//log.Println(sql)
}

func init() {
	_, err := loadTools()
	if err != nil {
		log.Fatalln("loadTools:", err)
	}
}

/**
 * 数据打码
 */
func (obj *Tool) Init() {
	if obj == nil {
		return
	}

	return
}

/**
 *  获取结构体，方便对结构体的数据进行初始化
 */
func GetTool() Tool {
	var obj Tool
	obj.DB = SqlDB.Conn
	fmt.Println(SqlDB.Database)
	return obj
}

// 设置User的表名为`profiles`
func (obj Tool) TableName() string {
	return "tools"
}

// 获取数据库链接
func (obj Tool) getDB() (db *gorm.DB) {
	db = obj.DB
	if db == nil {
		db = SqlDB.GetConn()
	}
	db = db.Table(obj.TableName())
	db = db.Select(obj.FieldList())
	return db
}

func (obj Tool) FieldList() (fieldList []string) {
	fieldList = SqlDB.GetFieldName(obj)
	return
}

// 传递过来的json数据初始化
func (obj *Tool) Insert() (id int64, err error) {
	err = obj.getDB().Create(&obj).Error
	if err == nil {
		//id = obj.ID
	}
	return
}

// 更新部分, 可以清空
func (obj *Tool) Update(id int64, updateData map[string]interface{}) (err error) {
	err = obj.Get(id)
	if err != nil {
		return
	}
	updateData = SqlDB.CleanChangeMap(updateData, obj)
	err = obj.getDB().Model(&obj).Updates(updateData).Error
	return
}

// 删除数据
func (obj *Tool) Delete(id int64) (err error) {
	err = obj.Get(id)
	if err != nil {
		return
	}
	err = obj.getDB().Delete(&obj).Error
	return
}

// 获取一个
func (obj *Tool) Get(id int64) (err error) {
	err = obj.getDB().First(&obj, id).Error
	if err == nil {
		obj.Init()
	}
	return
}

// 获取查询db
func (obj *Tool) QueryDB(db *gorm.DB, searchKey map[string]string) (returnDB *gorm.DB) {

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

func (obj *Tool) QueryOneByCondition(searchKey map[string]string) (err error) {
	db := obj.getDB()
	db = obj.QueryDB(db, searchKey)
	err = db.First(&obj).Error

	if err == nil {
		obj.Init()
	}

	return
}

// 列表
func (obj *Tool) List(searchKey map[string]string, page int64, pageSize int64) (results []Tool, err error) {
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

// 获取数量
func (obj *Tool) Count(searchKey map[string]string) (totalSize int64, err error) {
	db := obj.getDB().Where("deleted_at is null")
	err = obj.QueryDB(db, searchKey).Count(&totalSize).Error

	return
}
