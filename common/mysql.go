package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

/*
	数据源语法："用户名:密码@[连接方式](主机名:端口号)/数据库名"

	注意：open()在执行时不会真正的与数据库进行连接，只是设置连接数据库需要的参数
	ping()方法才是连接数据库
*/

// NewMysqlConn 创建 Msql 连接
func NewMysqlConn() (db *sql.DB, err error) {
	//"user:password@tcp(127.0.0.1:3306)/dbname"
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/imooc?charset=utf8")
	//if err != nil {
	//	panic(err)
	//}
	return
}

// GetResultRow 获取返回值，获取一条
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]string, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	for rows.Next() {
		//将行数据扫描到这个scanArgs中
		rows.Scan(scanArgs...)
		//scanArgs取了values地址，
		for i, v := range values {
			//fmt.Println("--------------GetResultRow-------------", v)
			if v != "" {
				//fmt.Println(reflect.TypeOf(col))
				record[columns[i]] = v
			}
		}
	}
	return record
}

// GetResultRows 获取所有
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}
