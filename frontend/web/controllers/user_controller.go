package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"imooc-Product/datamodels"
	"imooc-Product/encrypt"
	"imooc-Product/services"
	"imooc-Product/tool"
	"strconv"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session *sessions.Session
}

/*
// View 完成 `hero.Result` 接口。
// 它被用作替代返回值
// 包装模板文件名、布局、（任何）视图数据、状态码和错误。
// 它足够聪明地完成请求并将正确的响应发送给客户端。
*/

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		passWord = c.Ctx.FormValue("passWord")
	)
	//ozzo-validation
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: passWord,
	}

	_, err := c.Service.AddUser(user)
	c.Ctx.Application().Logger().Debug(err)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")
	return
}

func (c *UserController) GetLogin() mvc.Response {
	if c.Ctx.GetCookie("userid") != "" {
		return mvc.Response{
			Path: "/product/detail",
		}
	}
	return mvc.Response{
		Path: "user/login",
	}
}

/*
	Response
	它被用作替代返回值
	将状态码、内容类型、内容包装为字节或字符串
	和一个错误，它足够聪明地完成请求并将正确的响应发送给客户端。
*/

func (c *UserController) PostLogin() mvc.Response {
	//1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		passWord = c.Ctx.FormValue("passWord")
	)
	//2、验证账号密码正确
	user, isOk := c.Service.IsPwdSuccess(userName, passWord)
	//fmt.Println("----PostLogin !isOK-------", isOk, user)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	//3、写入用户ID到cookie中 引入加密 防止篡改
	tool.GlobalCookie(c.Ctx, "userid", strconv.FormatInt(user.ID, 10))
	//userid 转换为对应类型
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	//加密
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	//后面优化会去掉
	//c.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	//写入cookie 浏览器
	tool.GlobalCookie(c.Ctx, "sign", uidString)
	return mvc.Response{
		Path: "/product/detail",
	}

}
