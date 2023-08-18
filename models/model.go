package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	Name     string `orm:"unique"`
	Password string
	Email    string
	Articles []*Article `orm:"rel(m2m)"`
}

// 文章表和文章类型表是一对多
type Article struct {
	Id          int       `orm:"pk;auto"`
	Aname       string    `orm:"size(20)"`
	Atime       time.Time `orm:"auto_now"`
	Acount      int       `orm:"default(0);null"`
	Acontent    string
	Aimg        string
	ArticleType *ArticleType `orm:"rel(fk)"`
	Users       []*User      `orm:"reverse(many)"`
}
type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"` //设置一对多类型
}

func init() {
	orm.RegisterDataBase("default", "mysql", "root:1234@tcp(127.0.0.1:3306)/test1?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}
