package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"html/template"
	"imooc-Product/datamodels"
	"imooc-Product/rabbitmq"
	"imooc-Product/services"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	UserService    services.IUserService
	//消息队列
	RabbitMQ *rabbitmq.RabbitMQ
	Session  *sessions.Session
}

var (
	htmlOutPath  = "./web/htmlProductShow/" //生成的HTML保持目录
	templatePath = "./web/views/template/"  //静态文件模板目录
)

// GetGenerateHtml 生成模板
func (p *ProductController) GetGenerateHtml() {
	productIDString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productIDString)
	//1. 获取模板文件地址
	tmpPath := filepath.Join(templatePath, "product.html")
	contentTmp, err := template.ParseFiles(tmpPath)
	fmt.Println("2-------22-------", contentTmp)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2. 获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
	//3. 获取模板渲染数据
	product, err := p.ProductService.GetProductByID(int64(productID))
	fmt.Println(product, "---------222--22")
	//4. 生成静态文件
	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
}

//生成html静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	//1. 判断静态文件是否存在，存在就删除
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			fmt.Println("1--------------")
			ctx.Application().Logger().Debug(err)
		}
	}
	//2. 生成静态文件
	fmt.Println("2:-------------", fileName)
	//os.O_WRONLY 会利用html模板生成静态页面
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	defer file.Close()
	//将product 渲染到 静态文件file 中
	template.Execute(file, &product)
}

// Exist 判断文件是否存在，存在删除重新生成
func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// GetDetail 秒杀页面/商品详情页 需要用户登录才能看 curl //detail
func (p *ProductController) GetDetail() mvc.View {
	//固定产品
	product, err := p.ProductService.GetProductByID(4)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	fmt.Println("-----------产品路径------------", product)
	return mvc.View{
		//商品详情布局文件
		Layout: "shared/productLayout.html",
		//页面模板
		Name: "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

/*
	生产端代码
*/
// GetOrder func (p *ProductController) GetOrder() mvc.View {
//获取订单
func (p *ProductController) GetOrder() mvc.View {
	productIDString := p.Ctx.URLParam("productID")
	userIDString := p.Ctx.GetCookie("userid")
	//通过消息队列 缓解订单更新数据库压力
	//productIDString转换成int64
	productID, err := strconv.ParseInt(productIDString, 10, 64)
	fmt.Println("GetOrder:-------", productID)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//创建消息体
	message := datamodels.NewMessage(productID, userID)
	//类型转化 成rabbitmq可以传送的消息
	byteMessage, err := json.Marshal(message)
	fmt.Println()
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//使用Simple模式消息队列 实例定义好的模式
	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	user, err := p.UserService.GetUserByID(userID)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		//商品详情布局文件
		Layout: "shared/productLayout.html",
		//页面模板
		Name: "product/result.html",
		Data: iris.Map{
			"orderID":     user.NickName,
			"showMessage": "抢购成功！",
		},
	}
	//return byteMessage

	//productID由string转换成int
	//productID, err := strconv.Atoi(productIDString)
	//if err != nil {
	//	p.Ctx.Application().Logger().Debug(err)
	//}
	//product, err := p.ProductService.GetProductByID(int64(productID))
	//if err != nil {
	//	p.Ctx.Application().Logger().Debug(err)
	//}
	//var orderID int64
	//showMessage := "抢购失败！"
	////判断商品数量是否满足需求
	//if product.ProductNum > 0 {
	//	//扣除商品数量
	//	product.ProductNum -= 1
	//	//大流量下会出现超卖问题
	//	err := p.ProductService.UpdateProduct(product)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	}
	//	//创建订单
	//	userID, err := strconv.Atoi(userIDString)
	//	order := &datamodels.Order{
	//		UserID:     int64(userID),
	//		ProductID:  int64(productID),
	//		OderStatus: datamodels.OrderSuccess,
	//	}
	//
	//	//新建订单
	//	orderID, err = p.OrderService.InsertOrder(order)
	//	fmt.Println("------------order-----", orderID)
	//	if err != nil {
	//		p.Ctx.Application().Logger().Debug(err)
	//	} else {
	//		showMessage = "抢购购成功！"
	//	}
	//}
	//return mvc.View{
	//	//商品详情布局文件
	//	Layout: "shared/productLayout.html",
	//	//页面模板
	//	Name: "product/result.html",
	//	Data: iris.Map{
	//		"orderID":     ,
	//		"showMessage": ,
	//	},
	//}
}
