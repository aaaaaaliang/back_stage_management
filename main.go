package main

import (
	_ "aliangdemo/models"
	_ "aliangdemo/routers"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrePage", HandlePrePage)
	beego.AddFuncMap("ShowNextPage", HandleNextPage)
	beego.Run()
}

func HandlePrePage(data int) string {
	pageIndex := data - 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
func HandleNextPage(data int) string {
	pageIndex := data + 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
