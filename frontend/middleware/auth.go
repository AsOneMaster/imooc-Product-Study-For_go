package middleware

import "github.com/kataras/iris/v12"

func AuthConProduct(ctx iris.Context) {

	uid := ctx.GetCookie("userid")
	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录!")
		ctx.Redirect("/user/login")
		return
	}
	ctx.Application().Logger().Debug("已经登陆")
	//继续请求上下文
	ctx.Next()
}
