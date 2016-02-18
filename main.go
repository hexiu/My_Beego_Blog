package main

import (
	_ "beego/routers"
	"github.com/astaxie/beego"
	"beego/controllers"
	"beego/models"
	"github.com/astaxie/beego/orm"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug=true
	orm.RunSyncdb("default",false,true)
	beego.Router("/",&controllers.HomeController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.Router("/topic",&controllers.TopicController{})
	beego.Router("/category",&controllers.CategroyController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Run()

}


