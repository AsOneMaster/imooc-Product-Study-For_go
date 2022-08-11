package datamodels

// Order 订单属性
type Order struct {
	ID         int64 `sql:"ID"`
	UserID     int64 `sql:"userID"`
	ProductID  int64 `sql:"productID"`
	OderStatus int64 `sql:"oderStatus"`
}

//变量
const (
	OderWait     = iota //初始值为：0
	OrderSuccess        //1
	OrderFailed         //2
)
