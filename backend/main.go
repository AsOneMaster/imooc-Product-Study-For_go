package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
	"imooc-Product/backend/web/controllers"
	"imooc-Product/common"
	"imooc-Product/repositories"
	"imooc-Product/services"
)

/*
	秒杀系统后端管理项目启动文件
*/
func main() {
	//1. 创建iris 实例
	app := iris.New()
	//2. 设置错误模型
	app.Logger().SetLevel("debug")
	//3. 注册模板
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(
		true,
	)
	app.RegisterView(tmplate)
	//4. 设置模板静态文件目标 老版本 StaticWeb() 方法被移除
	app.HandleDir("/assets", "./web/assets")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(context iris.Context) {
		context.ViewData("message", context.Values().GetStringDefault("message", "访问页面出错"))
		context.ViewLayout("")
		context.View("shared/error.html")
	})

	//5. 连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//6. 注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//7. 启动服务
	app.Run(
		iris.Addr("localhost:9999"),
		//忽略iris框架的错误
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
