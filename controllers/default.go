
package controllers

import (
	"github.com/astaxie/beego"
	"beego/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
//	c.Data["Website"] = "jaxiu.cn"
//	c.Data["Email"] = "axiu@jaxiu.cn"
	c.Data["IsHome"]=true
	c.TplName = "home.html"
	c.Data["IsLogin"]=checkAccount(c.Ctx)

	topics,err:=models.GetAllTopics(c.Input().Get("cate"),true)
	if err!=nil {
		beego.Error(err.Error())
	} else {
		c.Data["Topics"]=topics
	}

	categories,err:=models.GetAllCategories()
	if err!=nil {
		beego.Error(err)
	}

	c.Data["Categories"]=categories

	/*

	c.Data["True"]=true
	c.Data["False"]=false
*/
/*
	type u struct {
		Name string
		Age int
		Sex string
	}

	user:=&u{
		Name :"hexiu",
		Age : 19,
		Sex :"Man",
	}
	c.Data["User"]=user

//	var nums []int =[1,2,3,4,5,6,7,8,9,0]
	nums:=[]int{1,2,3,4,5,6,7,8,9,0}
	c.Data["NUMS"]=nums
	c.Data["Hey"]="hi,guys!"
	c.Data["Html"]="<div>hello,beego!</div>"
*/
}
