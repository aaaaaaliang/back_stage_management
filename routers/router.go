package routers

import (
	"aliangdemo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	/*beego.InsertFilter("/Article/*", beego.BeforeRouter, FilterFunc)*/
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.MainController{})
	//注意：当实现了自定义的get请求方法，请求将不会访问默认方法
	beego.Router("/login", &controllers.MainController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/index", &controllers.MainController{}, "get:ShowIndex")
	beego.Router("/addArticle", &controllers.MainController{}, "get:ShowAdd;post:HandleAdd")
	beego.Router("/content", &controllers.MainController{}, "get:ShowContent")
	beego.Router("/update", &controllers.MainController{}, "get:ShowUpdate;post:HandleUpdate")
	beego.Router("/delete", &controllers.MainController{}, "get:HandleDelete")

	//添加类型
	beego.Router("/AddArticleType", &controllers.MainController{}, "get:ShowAddType;post:HandleAddType")
	//退出登陆
	beego.Router("/Logout", &controllers.MainController{}, "get:Logout")

}

/*var FilterFunc = func(ctx *context.Context) {
	username := ctx.Input.Session("username")
	if username == nil {
		ctx.Redirect(302, "/")
	}
}
*/
