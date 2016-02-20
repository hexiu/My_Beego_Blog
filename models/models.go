package models

import (
	"time"
	"os"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/Unknwon/com"

	"path"
	"fmt"
	"strconv"
	"strings"
)

const (
	_DB_NAME = "data/beeblog.db"
	_SQLITE3_DRIVER = "sqlite3"

)

type Category  struct {
	Id				int64
	Title			string
	Created 		time.Time	`orm:"index"`
	Views			int64		`orm:"index"`
	TopicTime 		time.Time   `orm:"index"`
	TopicCount		int64
	TopicLastUserId int64
}

type Topic struct {
	Id				int64
	Uid 			int64
	Title 			string
	Category        string
	Labels 			string
	Content 		string 		`orm:"size(5000)"`
	Attachment 		string
	Created			time.Time 	`orm:"index"`
	Updated 		time.Time  	`orm:"index"`
	Views 			int64 		`orm:"index"`
	Author  		string
	ReplyTime 		time.Time	`orm:"index"`
	ReplyCount 		int64
	ReplyLastUserId	int64
}

type Comment struct {
	Id 			int64
	Tid 		int64
	Name		string
	Content 	string			`orm:"size(1000)"`
	Created		time.Time 		`orm:"index"`
}

func RegisterDB()  {
	if(!com.IsExist(_DB_NAME)){
		os.MkdirAll(path.Dir(_DB_NAME),os.ModePerm)
		os.Create(_DB_NAME)
	}
	//注册模型
	orm.RegisterModel(new(Category),new(Topic),new(Comment))
	orm.RegisterDriver(_SQLITE3_DRIVER,orm.DRSqlite)
	orm.RegisterDataBase("default",_SQLITE3_DRIVER,_DB_NAME,10)
}

func AddCategory(name string) error {
	o:=orm.NewOrm()
	cate:=&Category{
		Title: name,
		Created : time.Now(),
		TopicTime : time.Now(),
	}
	qs:=o.QueryTable("category")
	err:= qs.Filter("title",name).One(cate)
	if err==nil {
		return err
	}

	fmt.Println(name,err)

	_,err =o.Insert(cate)

	fmt.Println(name,err)
	if err != nil {
		return err
	}

	return nil
}



func DelCategory(id string) error {
	cid,err:=strconv.ParseInt(id,10,64)
	if err != nil {
		return err
	}

	o:=orm.NewOrm()
	cate:=&Category{Id: cid}
	_,err = o.Delete(cate)

	return err
}

func GetAllCategories()  ([]*Category,error)  {
	o:=orm.NewOrm()
	cates := make([]*Category,0)

	qs:=o.QueryTable("category")

	_,err:=qs.All(&cates)

	return cates,err

}



func AddTopic(title,category,label,content,attachment string) error {
	//处理标签
	label="$"+strings.Join(strings.Split(label," "),"#$")+"#"
	//空格作为多个标签的分隔符
	// beego
	//$beego#
	//orm
	//$beego#$orm#

	// bee   beego


	o:=orm.NewOrm()

	topic:=&Topic{
		Title:title,
		Category:category,
		Labels:label,
		Content:content,
		Attachment:attachment,
		Created:time.Now(),
		Updated:time.Now(),
		ReplyTime:time.Now(),
	}

	_,err:=o.Insert(topic)
	if err != nil {
		return err
	}

	cate:=new(Category)
	qs:=o.QueryTable("category")
	err=qs.Filter("title",category).One(cate)
	if err == nil {
		//如果不存在，简单的忽略更新操作，存在更新
		cate.TopicCount++
		_,err=o.Update(cate)
	}

	return 	err

}

func GetAllTopics(cate,label string, isDesc bool) ([]*Topic,error)  {
	o:=orm.NewOrm()
	topics:=make([]*Topic,0)
	qs:=o.QueryTable("topic")

	var err error

	if isDesc {
		if len(cate)>0 {
			qs=qs.Filter("category",cate)
		}
		if len(label)>0 {
			qs=qs.Filter("labels__contains","$"+label+"#")
		}
		_,err=qs.OrderBy("-created").All(&topics)
	}else{
		_,err=qs.All(&topics)
	}

	return topics,err
}


func GetTopic(tid string) (*Topic,error) {
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil {
		return nil,err
	}

	o:=orm.NewOrm()
	topic:=new(Topic)

	qs:=o.QueryTable("topic")
	err=qs.Filter("id",tidNum).One(topic)
	if err!=nil{
		return nil,err
	}
	topic.Views++
	_,err=o.Update(topic)

	topic.Labels=strings.Replace(strings.Replace(topic.Labels,"#"," ",-1),"$","",-1)

	return topic,err
}


func ModifyTopic(tid,title,category,label,content,attachment string) error {

	label="$"+strings.Join(strings.Split(label," "),"#$")+"#"
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil {
		return err
	}

	var oldCate,oldAttach string

	o:=orm.NewOrm()
	topic:=&Topic{ Id: tidNum}
	if o.Read(topic) == nil {
		oldCate=topic.Category
		oldAttach=topic.Attachment
		topic.Title=title
		topic.Category=category
		topic.Labels=label
		topic.Content=content
		topic.Attachment=attachment
		topic.Updated=time.Now()
		_,err:=o.Update(topic)
		if err != nil {
			return err
		}
	}

	//删除附件
	if len(oldAttach)>0 {
		os.Remove(path.Join("attachment",oldAttach))
	}

	//更新分类统计
	if len(oldCate)>0 {
		cate:=new(Category)
		qs:=o.QueryTable("category")
		err=qs.Filter("title",oldCate).One(cate)
		if err == nil {
			//如果不存在，简单的忽略更新操作，存在更新
			cate.TopicCount--
			_,err=o.Update(cate)
		}
	}

	cate:=new(Category)
	qs:=o.QueryTable("category")
	err=qs.Filter("title",category).One(cate)
	if err == nil {
		//如果不存在，简单的忽略更新操作，存在更新
		cate.TopicCount++
		_,err=o.Update(cate)
	}

	return nil
}

func DeleteTopic(tid string) error {
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil {
		return err
	}

	var	oldCate string

	o:=orm.NewOrm()
	topic:=&Topic{Id: tidNum}
	if o.Read(topic) == nil {
		oldCate=topic.Category
		_,err:=o.Delete(topic)
		if err!=nil {
			return err
		}
		if len(oldCate)>0 {
			cate:=new(Category)
			qs:=o.QueryTable("category")
			err:=qs.Filter("title",oldCate).One(cate)
			if err == nil {
				cate.TopicCount--
				_,err=o.Update(cate)
			}
		}
	}

	return err
}


func AddReply(tid,nickname,content string) error {
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil {
		return err
	}
	reply:=&Comment{
		Tid:tidNum,
		Name:nickname,
		Content:content,
		Created:time.Now(),
	}

	o:=orm.NewOrm()
	_,err =o.Insert(reply)
	if err!=nil {
		return err
	}

	topic:=new(Topic)

	/*也可以这样实现
		topic := &Topic{Id:tidNum}
		if o.Read(topic)==nil {
			topic.ReplyCount++
			topic.ReplyTime=time.Now()
		}
		_,err=o.Update(topic)
		return err
	*/
	qs:=o.QueryTable("topic")
	err=qs.Filter("id",tidNum).One(topic)
	if err!=nil{
		return err
	}
	topic.ReplyCount++
	topic.ReplyTime=time.Now()
	_,err=o.Update(topic)
	return err
}

func GetAllReplies(tid string) (replies []*Comment ,err error) {
	tidNum,err:=strconv.ParseInt(tid,10,64)
	if err!=nil {
		return nil,err
	}

	replies = make([]*Comment,0)
	o:=orm.NewOrm()
	qs:=o.QueryTable("comment")
	_,err=qs.Filter("tid",tidNum).All(&replies)
	return replies,err
}

func DeleteReply(rid,tid string) error  {
	ridNum,err:=strconv.ParseInt(rid,10,64)
	if err!=nil {
		return err
	}

	var tidNum int64

	o:=orm.NewOrm()
	reply:=&Comment{Id:ridNum}

	if o.Read(reply)==nil {
		tidNum=reply.Tid
		_,err=o.Delete(reply)
		if err!=nil {
			return err
		}
	}
	replies := make([]*Comment,0)
	qs:=o.QueryTable("comment")
	_,err=qs.Filter("tid",tidNum).OrderBy("-Created").All(&replies)
	if err!=nil{
		return err
	}

	topic:= &Topic{Id:tidNum}
	if o.Read(topic) == nil {
		if int64(len(replies))==0 {
			topic.ReplyTime=time.Date(1970,1,1,0,0,0,0,time.UTC)
			topic.ReplyCount=0
			_,err=o.Update(topic)
			return err
		}
		topic.ReplyTime=replies[0].Created
		topic.ReplyCount=int64(len(replies))
		_,err=o.Update(topic)
	}

/*
	tidNum,err=strconv.ParseInt(tid,10,64)
	if err!=nil {
		return err
	}

	_,err=o.Delete(reply)

	fmt.Println(tidNum)

	topic:=new(Topic)
	qs:=o.QueryTable("topic")
	err=qs.Filter("id",tidNum).One(topic)
	if err!=nil{
		return err
	}

	topic.ReplyCount--
	topic.ReplyTime=time.Now()
	_,err=o.Update(topic)
*/
	return err

}