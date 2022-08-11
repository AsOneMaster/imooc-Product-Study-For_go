package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imooc-Product/common"
	"imooc-Product/datamodels"
	"imooc-Product/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
}

/*
	通过控制器方法的输入参数访问动态路径参数，不需要绑定。当你使用 iris 的默认语法来解析控制器处理程序时，你需要在方法后加上 “.” 字符，大写字母是一个新的子路径。 官网例子：
	 1  mvc.New(app.Party("/user")).Handle(new(user.Controller))
	 2
	 3 func(\*Controller) Get() - GET:/user.
	 4 func(\*Controller) Post() - POST:/user.
	 5 func(\*Controller) GetLogin() - GET:/user/login
	 6 func(\*Controller) PostLogin() - POST:/user/login
	 7 func(\*Controller) GetProfileFollowers() - GET:/user/profile/followers
	 8 func(\*Controller) PostProfileFollowers() - POST:/user/profile/followers
	 9 func(\*Controller) GetBy(id int64) - GET:/user/{param:long}
	10 func(\*Controller) PostBy(id int64) - POST:/user/{param:long}

*/
// curl GET: -i http://localhost:9999/product/all
func (p *ProductController) GetAll() mvc.View {
	productArray, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": productArray,
		},
	}
}

// PostUpdate 修改商品
// curl POST: -i http://localhost:9999/product/update
func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	//fmt.Println("--------PostUpdate1:--------------------", p.Ctx.Request().Form)
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//fmt.Println("--------PostUpdate2:--------------------", product)
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// GetAdd GET: curl -i http://localhost:9999/product/add
func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

// PostAdd POST: curl -i http://localhost:9999/product/add
func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// GetManager GET: curl -i http://localhost:9999/product/manager
func (p *ProductController) GetManager() mvc.View {
	//fmt.Println("--------GetManager:---------------", p.Ctx)
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// GetDelete GET: curl -i http://localhost:9999/product/delete
func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	isOk := p.ProductService.DeleteProductByID(id)
	if isOk {
		p.Ctx.Application().Logger().Debug("删除商品成功，ID为：" + idString)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败，ID为：" + idString)
	}
	p.Ctx.Redirect("/product/all")
}
