package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"imooc-Product/datamodels"
	"imooc-Product/services"
	"log"
	"sync"
)

// url 格式 amqp://账号（imoocuser）：密码（imoocuser）@rabbitmq服务器地址：端口号/vhost（imooc）

const MQURL = "amqp://imoocuser:imoocuser@127.0.0.1:5672/imooc" //为什么15672会报错？

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机
	Exchange string
	//key
	key string
	//连接信息
	MqUrl string
	//加锁 防抢占
	sync.Mutex
}

// NewRabbitMQ 创建RabbitMQ基础实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, key: key, MqUrl: MQURL}
	var err error
	//创建rabbitmq连接错误 [返回*connection类型]
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel() // rabbitmq.conn为*connection类型
	rabbitmq.failOnErr(err, "创建channel失败")
	return rabbitmq
}

// Destory 定义断开连接函数
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//定义错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple --------- step:1. 创建 ”  简单模式  “ 下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	//调用基础实例 simple模型用不到 exchange和key
	return NewRabbitMQ(queueName, "", "")
}

// PublishSimple ----------step:2. 简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message string) error {
	//1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建（保证队列存在，消息能发送到队列中）
	r.Lock() //保证channel安全
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		return err
	}
	//发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据exchange类型判断是否能找到符合routerkey规则的队列，如果无法找到符合条件的队列，那么会把发送的消息回退给发送者
		false,
		//如果为true，当exchange发送消息对队列，队列没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

// ConsumeSimple  ----------step:3. 简单模式下消费者
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	//1. 申请队列，如果队列不存在会自动创建，如果存在则跳过创建（保证队列存在，消息能发送到队列中）
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//消费者控流
	r.channel.Qos(
		1,     //当前消费者一次能接受最大消息数量1个
		0,     //服务器传递的最大容量（）以8位字节为单位
		false, //如果为true 对channel可用
	)

	//2. 接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		//可改手动修改【false】 控制速率
		false,
		//是否独有
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//是否阻塞 false为阻塞
		false,
		nil)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启用协程处理消息
	var num int
	go func() {
		for d := range msgs {

			//实现我们要处理的逻辑函数
			//log.Printf("Received a message:%s", d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			num++
			fmt.Println("Received a message:", string(d.Body), "第：", num)
			//数据库插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			//
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}
			//如果为true 表示确认所有未确认的消息
			//为false 表示确认当前消息
			d.Ack(false) //手动确认消息 正确才下一个
		}
	}()

	log.Printf("[*] Waiting for message,To exit press CTRL+C")
	<-forever
}
