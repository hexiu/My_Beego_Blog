package main

import (
	_ "beego/routers"
	"github.com/astaxie/beego"
	"beego/controllers"
	"beego/models"
	"github.com/astaxie/beego/orm"
	"os"
)

func init() {
	//注册数据库
	models.RegisterDB()
}

func main() {
	orm.Debug=true
	orm.RunSyncdb("default",false,true)

	//注册路由
	beego.Router("/",&controllers.HomeController{})
	beego.Router("/login",&controllers.LoginController{})
	beego.Router("/topic",&controllers.TopicController{})
	beego.Router("/category",&controllers.CategroyController{})
	beego.Router("/reply",&controllers.ReplyController{})
	beego.Router("/reply/add",&controllers.ReplyController{},"post:Add")
	beego.Router("/reply/delete",&controllers.ReplyController{},"get:Delete")
	beego.AutoRouter(&controllers.TopicController{})

	//创建附件目录
		os.Mkdir("attachment",os.ModePerm)

	//作为静态文件
//	beego.SetStaticPath("/attachment","attachment")
	//作为一个单独的控制器来处理
	beego.Router("/attachment/:all",&controllers.AttachController{})

	//启动Beego
	beego.Run()

}


