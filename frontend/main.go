package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc-Product/common"
	"imooc-Product/frontend/middleware"
	"imooc-Product/frontend/web/controllers"
	"imooc-Product/rabbitmq"
	"imooc-Product/repositories"
	"imooc-Product/services"
)

/*
	秒杀系统前端使用项目项目启动文件
*/
func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板
	app.HandleDir("/public", "./web/public")
	//访问生成好的html静态文件
	app.HandleDir("/html", "./web/htmlProductShow")
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
	}
	//sess := sessions.New(sessions.Config{
	//	Cookie:  "AdminCookie",
	//	Expires: 600 * time.Minute,
	//})
	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//注册db user对象
	user := repositories.NewUserManager("user", db)
	userService := services.NewUserService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	//注册rabbitmq实例
	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocProduct")

	//注册db product对象和order对象
	product := repositories.NewProductManager("product", db)
	order := repositories.NewOrderManager("order", db)

	//注册service
	productService := services.NewProductService(product)
	orderService := services.NewOrderService(order)

	//注册路由
	productPro := app.Party("/product")
	pro := mvc.New(productPro)

	//注册中间件
	productPro.Use(middleware.AuthConProduct)

	//注册绑定控制器 并注册服务与消息队列实例
	pro.Register(productService, orderService, rabbitmq)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
