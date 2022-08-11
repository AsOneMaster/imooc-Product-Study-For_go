package datamodels

// Message 简单的消息体
type Message struct {
	ProductID int64
	UserID    int64
}

// NewMessage 创建结构体
func NewMessage(productID, userID int64) *Message {
	return &Message{productID, userID}
}
