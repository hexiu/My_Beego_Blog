package controllers

import  (
	"github.com/astaxie/beego"
	"beego/models"
)

type TopicController  struct{
	beego.Controller
}


func (this *TopicController) Post ()  {
	if !checkAccount(this.Ctx) {
		this.Redirect("/login",302)
		return
	}

	title := this.Input().Get("title")
	category := this.Input().Get("category")
	content := this.Input().Get("content")
	tid := this.Input().Get("tid")
	var err error
	if len(tid)==0{
		err=models.AddTopic(title,category,content)

	}else{
		err=models.ModifyTopic(tid,title,category,content)
	}

	if err!=nil {
		beego.Error(err)
	}

	this.Redirect("/topic",302)

}



func (this *TopicController) Get () {
	this.Data["IsTopic"]=true
	this.Data["IsLogin"]=checkAccount(this.Ctx)
	this.TplName="topic.html"
	topics,err:=models.GetAllTopics("",false)
	if err!=nil {
		beego.Error(err.Error())
	} else {
		this.Data["Topics"]=topics
	}
}

func (this *TopicController) Add() {
	this.TplName="topic_add.html"
}

func (this *TopicController) View()  {
	this.TplName="topic_view.html"

	m:=this.Ctx.Input.Params()
	topic,err:=models.GetTopic(m["0"])
	if err!=nil{
		beego.Error(err)
		this.Redirect("/",302)
		return
	}

	this.Data["Topic"]=topic
	this.Data["Tid"]=m["0"]

	replies,err:=models.GetAllReplies(m["0"])
	if err != nil {
		beego.Error(err)
		return
	}

	this.Data["Replies"]=replies
	this.Data["IsLogin"]=checkAccount(this.Ctx)

}


func (this *TopicController) Modify()  {
	this.TplName="topic_modify.html"
	tid:=this.Input().Get("tid")
	topic,err:=models.GetTopic(tid)
	if err!=nil {
		beego.Error(err)
		this.Redirect("/",302)
		return
	}

	this.Data["Topic"]=topic
	this.Data["Tid"]=tid
}

func (this *TopicController) Delete(){
	if !checkAccount(this.Ctx){
		this.Redirect("/login",302)
		return
	}
	err:=models.DeleteTopic(this.Input().Get("tid"))
	if err!=nil{
		beego.Error(err)
	}
	this.Redirect("/",302)

}