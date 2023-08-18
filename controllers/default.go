package controllers

import (
	"aliangdemo/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/google/uuid"
	"math"
	"path"
	"strconv"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "login.html"
}

func (c *MainController) Post() {
	//数据处理
	//1.拿数据
	username := c.GetString("username")
	email := c.GetString("email")
	password := c.GetString("password")
	//2.对数据进行校验
	if username == "" || password == "" || email == "" {
		logs.Info("数据不能为空")
		//重定向
		c.Redirect("/register.html", 302)
		return
	}
	//3.插入数据库
	o := orm.NewOrm()
	//插入数据的结构体
	user := models.User{}
	user.Name = username
	user.Email = email
	user.Password = password
	_, err := o.Insert(&user)
	if err != nil {
		logs.Info("插入数据库失败")
		c.Redirect("/register", 302)
		return
	}
	//4.返回登陆，首页等操作
	c.TplName = "login.html"
}

func (c *MainController) ShowLogin() {
	name := c.Ctx.GetCookie("Name")
	c.Data["Name"] = name
	c.TplName = "login.html"
}

func (c *MainController) HandleLogin() {
	username := c.GetString("username")
	password := c.GetString("password")

	if username == "" || password == "" {
		logs.Info("输入数据不合法")
		c.TplName = "login.html"
	}
	o := orm.NewOrm()
	//赋值
	user := models.User{}
	user.Name = username
	err := o.Read(&user, "Name")
	if err != nil {
		logs.Info("查询失败")
		c.TplName = "login.html"
		return
	}
	c.SetSession("username", user.Name)
	c.Redirect("/index", 302)
}

// 显示首页内容
// controllers/default.go
func (c *MainController) ShowIndex() {
	username := c.GetSession("username")
	if username == nil {
		c.Redirect("/", 302) //如果没有session，回到登陆页面
		return
	}

	o := orm.NewOrm()
	var articles []models.Article
	qs := o.QueryTable("Article")
	/*qs.All(&articles)*/
	/*pageIndex1 := 1*/
	pageIndex := c.GetString("pageIndex")
	pageIndex1, err := strconv.Atoi(pageIndex)
	if err != nil {
		pageIndex1 = 1
	}

	//返回数据条目数
	count, err := qs.RelatedSel("ArticleType").Count()
	//获取总页数
	//一页有几条数据
	pageSize := 2
	start := pageSize * (pageIndex1 - 1)
	qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles) //1.pageSize 一页显示多少 2. start
	pageCount := float64(count) / float64(pageSize)
	pageCount1 := math.Ceil(pageCount)
	if err != nil {
		logs.Info("查询错误")
		return
	}
	FirstPage := false
	//首页 末页数据处理
	if pageIndex1 == 1 {
		FirstPage = true
	}
	LastPage := false
	if float64(pageIndex1) == pageCount1 {
		LastPage = true
	}
	//获取类型数据
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	c.Data["types"] = types

	//根据类型获取数据
	//接收数据
	typeName := c.GetString("select")
	/*logs.Info(typeName)*/
	//处理数据
	var articleswithtype []models.Article
	if typeName == "" {
		logs.Info("下拉框传递数据失败")
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articleswithtype)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType_TypeName", typeName).All(&articleswithtype)
	}
	//获取相应的文章列表

	logs.Info("count=", count)
	c.Data["LastPage"] = LastPage
	c.Data["FirstPage"] = FirstPage
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount1
	c.Data["articles"] = articleswithtype
	c.Data["pageIndex"] = pageIndex1
	c.TplName = "index.html"
}

// 显示添加文章页面

func (c *MainController) ShowContent() {
	username := c.GetSession("username")
	if username == nil {
		c.Redirect("/", 302) //如果没有session，回到登陆页面
		return
	}
	//获取文章ID
	id, err := c.GetInt("id")
	if err != nil {
		logs.Info("获取文章ID错误", err)
		return
	}
	//查询数据库
	o := orm.NewOrm()

	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		logs.Info("查询错误", err)
		return
	}
	//传递数据给视图
	c.Data["article"] = arti

	c.TplName = "content.html"
}

func (c *MainController) ShowAdd() {
	username := c.GetSession("username")
	if username == nil {
		c.Redirect("/", 302) //如果没有session，回到登陆页面
		return
	}
	//查询类型
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	c.Data["types"] = types

	c.TplName = "add.html"
}

//处理添加文章界面

func (c *MainController) HandleAdd() {
	//拿到数据
	articleName := c.GetString("articleName")
	articleContent := c.GetString("content")
	/*logs.Info(articleName, articleContent)*/
	f, h, err := c.GetFile("uploadname")
	defer f.Close()

	//要限定格式  限定大小  文件重新命名
	fileexit := path.Ext(h.Filename)
	uuid, err := uuid.NewRandom()
	if err != nil {
		logs.Info("UUID生成失败")
		return
	}
	uuidString := uuid.String()
	uniqueFilename := uuidString + fileexit
	filePath := "./static/img/" + uniqueFilename

	if fileexit != ".jpg" && fileexit != ".png" && fileexit != ".jpeg" {
		logs.Info("上传文件格式错误")
		return
	}
	if h.Size > 50000000 {
		logs.Info("上传文件过大")
		return
	}

	if err != nil {
		logs.Info("上传文件失败")
		return
	} else {
		c.SaveToFile("uploadname", "./static/img/"+filePath)
	}

	logs.Info(articleName, articleContent)
	//判断是否合法
	if articleContent == "" || articleName == "" {
		logs.Info("添加文章数据错误")
		return
	}
	//插入数据
	o := orm.NewOrm()
	arti := models.Article{}
	arti.Aname = articleName
	arti.Acontent = articleContent
	arti.Aimg = "/static/img/" + h.Filename

	//返回文章界面
	//给article对象赋值
	typeName := c.GetString("select")
	//类型判断
	if typeName == "" {
		logs.Info("下拉框数据错误")
		return
	}
	//获取type 对象
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err = o.Read(&artiType, "TypeName")
	if err != nil {
		logs.Info("类型获取失败")
		return
	}
	arti.ArticleType = &artiType

	_, err = o.Insert(&arti)
	if err != nil {
		logs.Info("插入数据库错误")
		return
	}
	c.Redirect("/index", 302)
}

// 显示编辑页面
func (c *MainController) ShowUpdate() {
	username := c.GetSession("username")
	if username == nil {
		c.Redirect("/", 302) //如果没有session，回到登陆页面
		return
	}
	//获取文章ID
	id, err := c.GetInt("id")
	if err != nil {
		logs.Info("获取文章ID错误", err)
		return
	}
	//查询数据库
	o := orm.NewOrm()

	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		logs.Info("查询错误", err)
		return
	}
	//传递数据给视图
	c.Data["article"] = arti

	c.TplName = "update.html"
}

//处理更新业务数据

func (c *MainController) HandleUpdate() {
	//拿到数据
	id, _ := c.GetInt("id")
	articleName := c.GetString("articleName")
	articleContent := c.GetString("content")
	/*logs.Info(articleName, articleContent)*/
	f, h, err := c.GetFile("uploadname")
	if err != nil {
		logs.Info("上传文件失败")
		return
	} else {
		defer f.Close()

		//要限定格式  限定大小  文件重新命名
		fileexit := path.Ext(h.Filename)
		uuid, err := uuid.NewRandom()
		if err != nil {
			logs.Info("UUID生成失败")
			return
		}
		uuidString := uuid.String()
		uniqueFilename := uuidString + fileexit
		filePath := "./static/img/" + uniqueFilename
		if fileexit != ".jpg" && fileexit != ".png" && fileexit != ".jpeg" {
			logs.Info("上传文件格式错误")
			return
		}
		if h.Size > 50000000 {
			logs.Info("上传文件过大")
			return
		}
		c.SaveToFile("uploadname", "./static/img/"+filePath)

		//对数据处理
		if articleName == "" || articleContent == "" {
			logs.Info("更新数据获取失败")
			return
		}
		//更新操作
		o := orm.NewOrm()

		arti := models.Article{Id: id}
		err = o.Read(&arti)
		if err != nil {
			logs.Info("查找数据错误")
			return
		}
		arti.Aname = articleName
		arti.Acontent = articleContent
		arti.Aimg = "./static/img/" + filePath

		_, err = o.Update(&arti, "Aname", "Acontent", "Aimg")
		if err != nil {
			logs.Info("更新数据错误")
			return
		}
		c.Redirect("/index", 302)
	}
}

// 删除操作
func (c *MainController) HandleDelete() {
	//拿到数据
	id, err := c.GetInt("id")
	if err != nil {
		logs.Info("获取ID错误")
		return
	}
	//执行删除操作
	o := orm.NewOrm()
	arti := models.Article{Id: id}
	err = o.Read(&arti)
	if err != nil {
		logs.Info("查询错误")
		return
	}
	o.Delete(&arti)
	//返回列表页面
	c.Redirect("/index", 302)
}

func (c *MainController) ShowAddType() {
	username := c.GetSession("username")
	if username == nil {
		c.Redirect("/", 302) //如果没有session，回到登陆页面
		return
	}
	//读取类型表，显示数据
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	//查询
	_, err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		logs.Info("查询类型错误")
	}
	c.Data["types"] = artiTypes
	c.TplName = "addType.html"
}

//处理添加类型业务

func (c *MainController) HandleAddType() {
	//获取数据
	typeName := c.GetString("typeName")
	//判断数据
	if typeName == "" {
		logs.Info("添加类型数据为空")
		return
	}
	//执行插入操作
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typeName
	_, err := o.Insert(&artiType)
	if err != nil {
		logs.Info("插入失败")
		return
	}
	//展示视图
	c.Redirect("/AddArticleType", 302)
}

// 退出登录
func (c *MainController) Logout() {
	//删除登陆状态
	c.DelSession("username")
	//退出登陆
	c.Redirect("/", 302)

}
