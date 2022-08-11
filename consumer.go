package main

import (
	"fmt"
	"imooc-Product/common"
	"imooc-Product/rabbitmq"
	"imooc-Product/repositories"
	"imooc-Product/services"
)

/*
	消费端代码
*/
func main() {
	//获取数据库实例
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	//创建product数据库操作实例
	product := repositories.NewProductManager("product", db)
	//创建 product service
	productService := services.NewProductService(product)
	//创建order 数据库操作实例
	order := repositories.NewOrderManager("order", db)
	//创建 order service
	orderService := services.NewOrderService(order)
	//消息队列消费者
	rabbitmqConsumer := rabbitmq.NewRabbitMQSimple("imoocProduct")
	//进行消费
	rabbitmqConsumer.ConsumeSimple(orderService, productService)

}
