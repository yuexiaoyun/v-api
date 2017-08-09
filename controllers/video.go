package controllers

import (
	/*"apiproject/models"
	"encoding/json"*/
	"apiproject/models"
	"context"
	"fmt"
	"github.com/astaxie/beego"
)

// Operations about object
type VideoController struct {
	beego.Controller
}

// @Title test
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /test [get]
func (this *VideoController) TestRouter() {
	this.Ctx.WriteString("测试路由")
	// vid := this.Input().Get("vid")

	// videoInfo := models.GetRawVideo(vid)
	// fmt.Println(videoInfo)
	/*this.Data["json"] = videoInfo
	this.ServeJSON()*/
	var vidStrList []string
	vidStrList = models.GetVidsByUid("50000336")
	fmt.Println(vidStrList)
}

// @Title test
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /shenjtlive [get]
func (this *VideoController) ShenJTLive() {
	versionCode := this.Input().Get("ver")

	yyuid := this.Input().Get("yyuid")
	limit := this.Input().Get("limit")
	page := this.Input().Get("page")

}

// @Title test
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /shenjtdetail [get]
func (this *VideoController) ShenJTLive() {
	versionCode := this.Input().Get("ver")

	vid := this.Input().Get("vid")

	models.GetByVid(vid)

}
