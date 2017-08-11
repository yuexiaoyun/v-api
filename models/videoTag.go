package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(VideoTag))
}

type VideoTag struct {
	Id         int64 `orm:"column(id)"`
	Vid        int64 `orm:"column(vid)"`
	Tag        string `orm:"column(tag)"`
	Priority   int `orm:"column(priority)"`
	CreateTime int64 `orm:"column(create_time)"`
	UpdateTime int64  `orm:"column(update_time)"`
}

func (t *VideoTag) TableName() string {
	return "v_huya_video_tags"
}

func GetTagByVid(vid int64) ([]VideoTag) {
	var videoTag []VideoTag
	o := orm.NewOrm()
	qs := o.QueryTable("v_huya_video_tags")
	qs.Filter("vid", vid)
	qs.All(&videoTag)
	return videoTag

}
