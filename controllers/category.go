package controllers

import (
	"github.com/astaxie/beego"
	"beego/models"
)


type CategroyController  struct {
	beego.Controller
}

func (this *CategroyController) Get(){
	op:=this.Input().Get("op")

	switch op {
		case "add" :
			name:=this.Input().Get("name")
			if len(name) == 0 {
				break
			}
			err:=models.AddCategory(name)

			if err != nil {
				beego.Error(err)
			}

			this.Redirect("/category",301)
			return

		case "del" :
			id:=this.Input().Get("id")
			if len(id) == 0 {
				break
			}

			err:=models.DelCategory(id)
			if err!=nil {
				beego.Error(err)
			}
			this.Redirect("/category",301)
			return
	}

	this.Data["IsCategory"]=true
	this.TplName="category.html"
	this.Data["IsLogin"] = checkAccount(this.Ctx)

	var err error
	this.Data["Categories"],err = models.GetAllCategories()
	if err!=nil {
		beego.Error(err)
	}
}


