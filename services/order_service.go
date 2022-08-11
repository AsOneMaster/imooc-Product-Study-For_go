package services

import (
	"imooc-Product/datamodels"
	"imooc-Product/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	DeleteOrderByID(int64) bool
	InsertOrder(order *datamodels.Order) (int64, error)
	// InsertOrderByMessage 消息队列 增加订单
	InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error)
	UpdateOrder(order *datamodels.Order) error
}

type OrderService struct {
	OrderRepository repositories.IOrder
}

// NewOrderService 初始化Order服务函数
func NewOrderService(repository repositories.IOrder) IOrderService {
	return &OrderService{OrderRepository: repository}
}

// GetOrderByID  通过ID查询订单
func (o *OrderService) GetOrderByID(orderID int64) (order *datamodels.Order, err error) {
	order, err = o.OrderRepository.SelectByKey(orderID)
	return
}

// GetAllOrder  查询所有订单
func (o *OrderService) GetAllOrder() (orders []*datamodels.Order, err error) {
	orders, err = o.OrderRepository.SelectAll()
	return
}

//GetAllOrderInfo 查询与订单有关的货物信息
func (o *OrderService) GetAllOrderInfo() (orderMap map[int]map[string]string, err error) {
	orderMap, err = o.OrderRepository.SelectAllWithInfo()
	return
}

// DeleteOrderByID  订单删除
func (o *OrderService) DeleteOrderByID(orderID int64) (isDelete bool) {
	isDelete = o.OrderRepository.Delete(orderID)
	return
}

// InsertOrder  订单添加
func (o *OrderService) InsertOrder(order *datamodels.Order) (orderID int64, err error) {
	orderID, err = o.OrderRepository.Insert(order)
	return
}

// InsertOrderByMessage 通过消息队列增加订单
func (o *OrderService) InsertOrderByMessage(message *datamodels.Message) (orderID int64, err error) {
	order := &datamodels.Order{
		UserID:     message.UserID,
		ProductID:  message.ProductID,
		OderStatus: datamodels.OrderSuccess,
	}
	orderID, err = o.InsertOrder(order)
	//fmt.Println("orderID生成：", orderID)
	return
}

// UpdateOrder  订单更新
func (o *OrderService) UpdateOrder(order *datamodels.Order) (err error) {
	err = o.OrderRepository.Update(order)
	return
}
