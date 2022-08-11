package repositories

import (
	"database/sql"
	"imooc-Product/common"
	"imooc-Product/datamodels"
	"strconv"
)

//1. 先开发接口
//2. 实现接口

// IProduct 定义接口规范
type IProduct interface {
	// Conn 连接数据库
	Conn() error
	Insert(product *datamodels.Product) (int64, error)
	Delete(id int64) bool
	Update(product *datamodels.Product) error
	SelectByKey(id int64) (product *datamodels.Product, err error)
	SelectAll() ([]*datamodels.Product, error)
	// SubProductNum 引入消息队列增加的 扣除商品数量的方法
	SubProductNum(productID int64) error
}

// ProductManager 实现接口的 结构体
type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

// NewProductManager 创建构造函数 返回类型为接口类型
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

/*
	以下是ProductManager实现接口的所有方法
*/
// Conn 数据库连接------------------------------
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// Insert 商品添加
func (p *ProductManager) Insert(product *datamodels.Product) (productID int64, err error) {
	//1. 判断数据库连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	//2. 准备sql
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	//获取预处理语句对象
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	//3. 调用预处理语句 传入参数
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return
	}
	// LastInsertId 返回数据库为响应命令而生成的整数。 通常，这将来自插入新行时的“自动增量”列。
	productID, err = result.LastInsertId()
	return

}

// Delete 商品删除
func (p *ProductManager) Delete(productID int64) bool {
	//1. 判断数据库连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM product WHERE ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(productID)
	if err != nil {
		return false
	}
	return true
}

// Update 商品更新
func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	//1. 判断数据库连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	sql := "UPDATE product SET productName=?,productNum=?,productImage=?,productUrl=? WHERE ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ID)
	if err != nil {
		return
	}
	return nil
}

// SelectByKey 通过ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (product *datamodels.Product, err error) {
	//1. 判断数据库连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	//strconv.FormatInt(productID, 10)转换为10进制
	sql := "SELECT ID,productName,productNum,productImage,productUrl FROM " + p.table + " WHERE ID=" + strconv.FormatInt(productID, 10)
	row, err := p.mysqlConn.Query(sql)
	if err != nil {
		return
	}
	//GetResultRow转换函数
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	product = &datamodels.Product{}
	common.DataToStructByTagSql(result, product)
	return
}

// SelectAll 查询所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, err error) {
	//1. 判断数据库连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	sql := "SELECT ID,productName,productNum,productImage,productUrl FROM " + p.table
	rows, err := p.mysqlConn.Query(sql)
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
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}

// SubProductNum 扣除商品数量
func (p *ProductManager) SubProductNum(productID int64) (err error) {
	if err = p.Conn(); err != nil {
		return
	}
	sql := "update " + p.table + " set " + " productNum=productNum-1 where ID=" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return
	}
	_, err = stmt.Exec()
	return
}
