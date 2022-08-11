package repositories

import (
	"database/sql"
	"fmt"
	"imooc-Product/common"
	"imooc-Product/datamodels"
	"strconv"
)

type IOrder interface {
	Conn() error
	Insert(order *datamodels.Order) (int64, error)
	Delete(int642 int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(int642 int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, db *sql.DB) IOrder {
	return &OrderManager{table: table, mysqlConn: db}
}

/*
	以下是OderManager实现接口的所有方法
*/
// Conn 数据库连接------------------------------
func (o *OrderManager) Conn() (err error) {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return
}

// Insert 订单添加
func (o *OrderManager) Insert(order *datamodels.Order) (orderID int64, err error) {
	//1. 判断数据库连接是否存在
	if err = o.Conn(); err != nil {
		return
	}
	//2. 准备sql
	sql := "INSERT `order` SET userID=?,productID=?,orderStatus=?"
	//获取预处理语句对象
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	//3. 调用预处理语句 传入参数
	fmt.Println()
	result, err := stmt.Exec(order.UserID, order.ProductID, order.OderStatus)
	if err != nil {
		return
	}
	// LastInsertId 返回数据库为响应命令而生成的整数。 通常，这将来自插入新行时的“自动增量”列。
	orderID, err = result.LastInsertId()
	return

}

// Delete 订单删除
func (o *OrderManager) Delete(orderID int64) bool {
	//1. 判断数据库连接是否存在
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM order WHERE ID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(orderID)
	if err != nil {
		return false
	}
	return true
}

// Update 订单更新
func (o *OrderManager) Update(order *datamodels.Order) (err error) {
	//1. 判断数据库连接是否存在
	if err = o.Conn(); err != nil {
		return
	}
	sql := "UPDATE order SET userID=?,productID=?,oderStatus=? WHERE ID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}

	_, err = stmt.Exec(order.UserID, order.ProductID, order.OderStatus)
	if err != nil {
		return
	}
	return nil
}

// SelectByKey 通过ID查询订单
func (o *OrderManager) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	//1. 判断数据库连接是否存在
	if err = o.Conn(); err != nil {
		return
	}
	//strconv.FormatInt(productID, 10)转换为10进制
	sql := "SELECT ID,userID,productID,oderStatus FROM " + o.table + " WHERE ID=" + strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	if err != nil {
		return
	}
	//GetResultRow转换函数
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

// SelectAll 查询所有订单
func (o *OrderManager) SelectAll() (orderArray []*datamodels.Order, err error) {
	//1. 判断数据库连接是否存在
	if err = o.Conn(); err != nil {
		return
	}
	sql := "SELECT ID,userID,productID,oderStatus FROM " + o.table
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return
	}
	//Close 关闭行，防止进一步枚举。 如果调用 Next 并返回 false 并且没有进一步的结果集，则 Rows 将自动关闭，并足以检查 Err 的结果。 Close 是幂等的，不影响 Err 的结果
	defer rows.Close()
	if err != nil {
		return
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

// SelectAllWithInfo 查询与订单有关联的货物信息
func (o *OrderManager) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	//1. 判断数据库连接是否存在
	if err = o.Conn(); err != nil {
		return
	}
	//imooc.order和order
	sql := "Select o.ID,o.userID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return
	}

	orderMap = common.GetResultRows(rows)
	return
}
