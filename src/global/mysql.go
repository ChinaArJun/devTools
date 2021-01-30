package global

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	gormlog "github.com/onrik/logrus/gorm"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"oktools/src/conf"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SqlDB SqlConn

type SqlConn struct {
	Database string
	Conn     *gorm.DB

	//这里是可增加的固定变量, 请注意不可无限量增大, 防止内存溢出
	dataFieldList   sync.Map
	dataFieldStruct sync.Map
	dataTableList   sync.Map
}

func init() {
	SqlDB.Connect()
}

func (obj *SqlConn)Connect() {
	dbConf := conf.Conf.Mysqldb

	var err error
	//var sqlHostString string
	sqlStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConf.Username,
		dbConf.Password, dbConf.Host, dbConf.Port, dbConf.DbName)
	//sqlHostString = dbConf.Username + ":" + dbConf.Password + "@" + "(" + dbConf.Host + ":" + dbConf.Port + ")"+ "/" +
	//	dbConf.DbName += "?charset=utf8&parseTime=True&loc=Local&timeout=10ms"
	log.Println(sqlStr)
	obj.Database = dbConf.DbName
	obj.Conn, err = gorm.Open("mysql", sqlStr)
	if err != nil {
		fmt.Println(err.Error())
		panic("数据库启动失败:" + err.Error())
	}
	intType, err := strconv.Atoi(dbConf.MaxId)
	if err != nil {
		obj.Conn.DB().SetMaxIdleConns(intType)
	}
	openConns, err := strconv.Atoi(dbConf.MaxOpen)
	if err != nil {
		obj.Conn.DB().SetMaxOpenConns(openConns)
	}

	if err != nil {
		fmt.Printf("mysql connect error %v\n", err)
	}

	if obj.Conn.Error != nil {
		fmt.Printf("\ndatabase error %v", obj.Conn.Error)
	}

	// 表名禁用复数
	obj.Conn.SingularTable(true)

	obj.Conn.SetLogger(gormlog.New(log.WithFields(log.Fields{"app": "gorm"})))
	obj.Conn.LogMode(true)

	log.Print("mysql connect success!")

	//db, err := gorm.Open("mysql", "root:amo1994@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	//if err!= nil{
	//	panic(err)
	//}

	//rows,err := obj.Conn.DB().Query("select * from python_fund")

	//for rows.Next(){
	//	//row.Scan(...)
	//}
	//rows.Close()
}
//func (obj *SqlConn) Close() {
//	obj.Conn.Close()
//}

func (obj *SqlConn) GetConn() *gorm.DB {
	return obj.Conn
}

// 获取结构体的字段
func (obj *SqlConn) GetFieldName(structName interface{}) (result []string) {

	t := reflect.TypeOf(structName)
	if v, ok := obj.dataFieldList.Load(t.String()); ok {
		return v.([]string)
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}

	ObjType := make(map[string]interface{})
	j, _ := json.Marshal(structName)
	json.Unmarshal(j, &ObjType)

	noInString := "_show"
	for k, _ := range ObjType {
		if len(k) >= len(noInString) && k[len(k)-len(noInString):] == noInString {
			continue
		}
		result = append(result, k)
	}

	if _, ok := t.FieldByName("CreatedAt"); ok {
		result = append(result, "created_at")
	}
	if _, ok := t.FieldByName("UpdatedAt"); ok {
		result = append(result, "updated_at")
	}

	if _, ok := t.FieldByName("DeletedAt"); ok {
		result = append(result, "deleted_at")
	}

	obj.dataFieldList.Store(t.String(), result)
	return
}

/*
// 获取结构图的字段及类型
//	字段不允许命名 _show 以 "_show" 结尾的数据属于业务数据
*/
func (obj *SqlConn) GetFieldStruct(structName interface{}) (result map[string]string) {

	result = make(map[string]string)

	t := reflect.TypeOf(structName)
	if v, ok := obj.dataFieldStruct.Load(t.String()); ok {
		return v.(map[string]string)
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}

	ObjType := make(map[string]interface{})
	j, _ := json.Marshal(structName)
	json.Unmarshal(j, &ObjType)

	for k, v := range ObjType {
		if v == nil {
			continue
		}
		result[k] = fmt.Sprintf("%T", v)
	}

	obj.dataFieldStruct.Store(t.String(), result)
	return
}

func (obj SqlConn) viewFieldString(fieldPrefix string, fieldMap map[string]string) string {
	sqlString := ""
	for k, v := range fieldMap {
		if sqlString != "" {
			sqlString += ","
		}

		sqlString += fieldPrefix
		sqlString += k
		sqlString += " as "
		sqlString += v
	}
	return sqlString
}

// 清理更新数据里面指定敏感字段
func (obj SqlConn) CleanChangeMap(updateMap map[string]interface{}, data interface{}) (returnMap map[string]interface{}) {

	delete(updateMap, "id")
	delete(updateMap, "created_at")
	delete(updateMap, "updated_at")
	delete(updateMap, "deleted_at")

	// 根据数据库的格式转换下更新的格式
	structMap := obj.GetFieldStruct(data)

	/*
		for k,v := range updateMap{
			fmt.Println( k, " = ", v, fmt.Sprintf("%T",v));
		}
	*/

	if len(updateMap) < len(structMap) {
		for k, sType := range structMap {
			if value, ok := updateMap[k]; ok {
				cType := fmt.Sprintf("%T", value)
				if sType != cType {
					//json过来的格式与数据库格式不一致
					switch sType {
					case "float64":
						if fmt.Sprintf("%T", value) == "string" {
							if tmp, err := strconv.Atoi(value.(string)); err == nil {
								updateMap[k] = int64(tmp)
							}
						}
						break
					case "string":
						if fmt.Sprintf("%T", value) == "float64" {
							updateMap[k] = strconv.Itoa(int(value.(float64)))
						}
						break
					}

				}
			}
		}
	} else {
		for k, value := range updateMap {
			if sType, ok := updateMap[k]; ok {
				cType := fmt.Sprintf("%T", value)
				if sType != cType {
					//json过来的格式与数据库格式不一致
					switch sType {
					case "float64":
						if fmt.Sprintf("%T", value) == "string" {
							if tmp, err := strconv.Atoi(value.(string)); err == nil {
								updateMap[k] = int64(tmp)
							}
						}
						break
					case "string":
						if fmt.Sprintf("%T", value) == "float64" {
							updateMap[k] = strconv.Itoa(int(value.(float64)))
						}
						break
					}

				}
			}
		}
	}

	/*
		for k,v := range updateMap{
			fmt.Println( k, " = ", v, fmt.Sprintf("%T",v));
		}
		panic("" )
	*/
	returnMap = updateMap
	return
}

/*
	// 初始化DB的查询条件, 根据map
	// 查询扩展
	字段中如果字段 为 realname	则默认支持以下查询
	realname 			realname	=	value
	realname_min 		realname	>= 	value
	realname_max 		realname	<= 	value
	realname_is_null	realname is null
	realname_is_empty	realname 	= 	''
	realname_not		realname	<> 	value
	realname_like		realname	like '%value%'
	realname_left_like	realname	like '%value'
	realname_right_like	realname	like 'value%'


	order
		name asc|desc
		FIELD name,0,1,2,-1,4,asc|desc
*/
//	如果不想用这个, 请直接在类中写, 这样可以增加效率或执行速度
func (obj SqlConn) QueryDbInit(db *gorm.DB, FieldMap []string, searchkey map[string]string) (returnDb *gorm.DB) {

	var tmpFieldMap = make(map[string]int64)
	for _, key := range FieldMap {
		key = strings.Trim(key, " ")
		if key == "" {
			continue
		}
		tmpFieldMap[key] = 1

		if value, ok := searchkey[key]; ok && value != "" {
			db = db.Where(key+" = ?", value)
		}
		//做in查询时，通过，分割
		if value, ok := searchkey[key+"_in"]; ok && value != "" {
			vals := strings.Split(value, ",")
			db = db.Where(key+" in (?)", vals)
		}
		if value, ok := searchkey[key+"_in_not"]; ok && value != "" {
			vals := strings.Split(value, ",")
			db = db.Where(key+" not in (?)", vals)
		}
		if value, ok := searchkey[key+"_min"]; ok && value != "" {
			db = db.Where(key+" >= ?", value)
		}
		if value, ok := searchkey[key+"_max"]; ok && value != "" {
			db = db.Where(key+" <= ?", value)
		}

		if value, ok := searchkey[key+"_gt"]; ok && value != "" {
			db = db.Where(key+" > ?", value)
		}
		if value, ok := searchkey[key+"_lt"]; ok && value != "" {
			db = db.Where(key+" < ?", value)
		}

		if value, ok := searchkey[key+"_is_null"]; ok && value != "" {
			db = db.Where(key + " is null")
		}

		if value, ok := searchkey[key+"_is_not_null"]; ok && value != "" {
			db = db.Where(key + " is not null")
		}

		if value, ok := searchkey[key+"_is_empty"]; ok && value != "" {
			db = db.Where(key + " = ''")
		}

		if value, ok := searchkey[key+"_is_not_empty"]; ok && value != "" {
			db = db.Where(key + " != ''")
		}

		if value, ok := searchkey[key+"_or"]; ok && value != "" {
			db = db.Or(key+" = ?", value)
		}

		// 全文检索 尽量少用
		if value, ok := searchkey[key+"_not"]; ok && value != "" {
			db = db.Where(key+" <> ?", value)
		}
		if value, ok := searchkey[key+"_like"]; ok && value != "" {
			db = db.Where(key+" LIKE ?", "%"+value+"%")
		}
		if value, ok := searchkey[key+"_left_like"]; ok && value != "" {
			db = db.Where(key+" LIKE ?", value+"%")
		}
		if value, ok := searchkey[key+"_right_like"]; ok && value != "" {
			db = db.Where(key+" LIKE ?", "%"+value)
		}

	}

	if value, ok := searchkey["order_by"]; ok && value != "" {
		value = strings.Trim(value, " ")
		chrstr := strings.Split(value, ",")
		for _, ov := range chrstr {
			ov = strings.Trim(ov, " ")
			if ov == "" {
				continue
			}
			order_array := strings.Split(ov, " ")

			if len(order_array) <= 2 {
				order_by_key := order_array[0]
				if order_by_key == "" {
					continue
				}
				//修改防止按order by status=5某个值排序的被过滤
				if _, ok := tmpFieldMap[order_by_key]; !ok && !strings.Contains(order_by_key, "=") {
					continue
				}

				order_by_value := "asc"
				if len(order_array) == 2 {
					order_by_value = strings.ToLower(order_array[1])
				}
				if order_by_value != "asc" && order_by_value != "desc" {
					continue
				}

				db = db.Order(order_by_key + " " + order_by_value)
			}

			if len(order_array) > 2 {
				if "FIELD" == strings.ToUpper(order_array[0]) {
					order_by_key := order_array[1]
					if order_by_key == "" {
						continue
					}
					if _, ok := tmpFieldMap[order_by_key]; !ok {
						continue
					}

					sqlString := ""
					for k, v := range order_array {
						//为了安全,只支持数字
						if k >= 2 && v != "" {
							if _, err := strconv.Atoi(v); err == nil {
								sqlString += "," + v
							}
						}
					}

					thisOrderDesc := strings.ToLower(order_array[len(order_array)-1])
					if thisOrderDesc != "asc" && thisOrderDesc != "desc" {
						thisOrderDesc = " asc"
					}

					if sqlString != "" {
						sqlString = "FIELD(" + order_by_key + sqlString + ") " + thisOrderDesc
						db = db.Order(sqlString)
					}
				}
			}
		}
	}

	returnDb = db
	return
}

/**
 *
 * 备份MySql数据库
 * @param 	host: 			数据库地址: localhost
 * @param 	port:			端口: 3306
 * @param 	user:			用户名: root
 * @param 	password:		密码: root
 * @param 	databaseName:	需要备份的数据库名数组 []string
 * @param 	sqlPath:		备份SQL存储路径: D:/backup/test/
 * @return 	backupPath
 *
 */
func (obj SqlConn) BackupMySql(host string, port string, user string, password string, databaseName []string, sqlPath string) (err error, filePath string) {
	// 获取根目录
	rootDir, _ := filepath.Abs(`.`)
	if sqlPath[len(sqlPath)-1:] != "/" {
		sqlPath += "/"
	}
	if rootDir[len(rootDir)-1:] != "/" {
		rootDir += "/"
	}

	pathNameTmp := time.Now().Format("20060102")
	pathName := rootDir + pathNameTmp + "/"

	//创建备份文件夹
	os.Mkdir(pathName, os.ModePerm)

	for _, name := range databaseName {
		obj.BackupMySqlDb(host, port, user, password, name, "", pathName)
	}

	fileName := pathNameTmp + ".tar.gz"

	//压缩备份文件
	cmd := exec.Command("tar", "-zcvf", fileName, pathNameTmp)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	_, err = ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	os.RemoveAll(pathName)
	os.Rename(fileName, sqlPath+fileName)

	return
}

/**
 *
 * 备份MySql数据库
 * @param 	host: 			数据库地址: localhost
 * @param 	port:			端口: 3306
 * @param 	user:			用户名: root
 * @param 	password:		密码: root
 * @param 	databaseName:	需要被分的数据库名: test
 * @param 	tableName:		需要备份的表名: user
 * @param 	sqlPath:		备份SQL存储路径: D:/backup/test/
 * @return 	backupPath
 *
 */
func (obj SqlConn) BackupMySqlDb(host, port, user, password, databaseName, tableName, sqlPath string) (err error, filePath string) {

	var cmd *exec.Cmd

	if tableName == "" {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName)
	} else {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName, tableName)
	}

	stdout, err := cmd.StdoutPipe()
	defer stdout.Close() // 保证关闭输出流
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	var backupPath string
	if tableName == "" {
		backupPath = sqlPath + databaseName + ".sql"
	} else {
		backupPath = sqlPath + databaseName + "_" + tableName + ".sql"
	}
	err = ioutil.WriteFile(backupPath, bytes, 0644)
	return
}

// 根据sqlstring 获取数据
func (obj *SqlConn) QueryData(sqlString string) (returnMap []map[string]string, err error) {
	rows, err := obj.Conn.Raw(sqlString).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]string)
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})

			if fmt.Sprintf("%T", *val) == "[]uint8" {
				m[colName] = string((*val).([]uint8))
			}

			if fmt.Sprintf("%T", *val) == "int64" {
				m[colName] = strconv.Itoa(int((*val).(int64)))
			}

		}
		returnMap = append(returnMap, m)
		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
	}
	return
}
