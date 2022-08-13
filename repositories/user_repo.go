package repositories

import (
	"database/sql"
	"github.com/kataras/iris/v12/core/errgroup"
	"imooc-Product/common"
	"imooc-Product/datamodels"
	"strconv"
)

type IUser interface {
	Conn() error
	Get(userID int64) (user *datamodels.User, err error)
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

type UserManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserManager(table string, db *sql.DB) IUser {
	return &UserManager{table, db}
}

func (u *UserManager) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManager) Get(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	sql := "Select * from " + u.table + " where ID=?"
	//fmt.Println("---------------selectName1:------------", userID, sql)
	rows, errRows := u.mysqlConn.Query(sql, userID)
	//str, _ := rows.Columns()
	//fmt.Println("---------------selectName2:------------", str)
	//sql.Rows 类型用完了要关闭
	defer rows.Close()
	if errRows != nil {
		return &datamodels.User{}, errRows
	}

	result := common.GetResultRow(rows)

	if len(result) == 0 {
		return &datamodels.User{}, errgroup.New("用户不存在！")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)

	return
}

func (u *UserManager) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errgroup.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	sql := "Select * from " + u.table + " where userName=?"
	//fmt.Println("---------------selectName1:------------", userName, sql)
	rows, errRows := u.mysqlConn.Query(sql, userName)
	//str, _ := rows.Columns()
	//fmt.Println("---------------selectName2:------------", str)
	//sql.Rows 类型用完了要关闭
	defer rows.Close()
	if errRows != nil {
		return &datamodels.User{}, errRows
	}

	result := common.GetResultRow(rows)

	if len(result) == 0 {
		return &datamodels.User{}, errgroup.New("用户不存在！")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)

	return
}

func (u *UserManager) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	sql := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return userId, errStmt
	}
	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errResult != nil {
		return userId, errResult
	}
	return result.LastInsertId()
}

func (u *UserManager) SelectByID(userId int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errgroup.New("用户不存在！")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
