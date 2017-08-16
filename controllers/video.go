package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"v-api/models"
	"encoding/json"
	"fmt"
)

// Operations about Video
type VideoController struct {
	beego.Controller
}

// @Title TestRouter
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /test [get]
func (this *VideoController) TestRouter() {
	//this.Ctx.WriteString("测试路由")
	yyuid := this.Input().Get("uid");
	vidStrList := models.GetVidsByUid(yyuid)
	fmt.Println(vidStrList)
	this.Data["json"] = "ok"
	this.ServeJSON()
}

// @Title ShenJTLive
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /shenjtlive [get]
func (this *VideoController) ShenJTLive() {
	versionCode := this.Input().Get("ver")


	yyuid := this.Input().Get("yyuid")
	limit := this.Input().Get("limit")
	page := this.Input().Get("page")

	limitInt,limitErr := strconv.ParseInt(limit,10,64)
	if limitErr != nil {
		beego.Error("limit is not a number");
	}
	pageInt,pageErr := strconv.ParseInt(page,10,64)
	if pageErr != nil {
		beego.Error("page is not a number");
	}
	if len(yyuid) == 0 || limitInt < 1 || pageInt < 1 {
		this.Data["json"] = "no data"
		this.ServeJSON()
	} else{
		versionCodeInt,versionErr := strconv.ParseInt(versionCode,10,64)
		if versionErr != nil {
			beego.Error("versioncode is not a number");
		}
		if versionCodeInt == 2 {

		}
	}

}


func shenJTLiveV2(){

}

// @Title ShenJTDetail
// @Description 获取视频详情
// @Param   vid     query   int false       "vid"
// @Success 200 {string} json success-视频详情
// @router /shenjtdetail [get]
func (this *VideoController) ShenJTDetail() {
	cacheKye := "video_detail_"
	vid := this.Input().Get("vid")
	if vid != ""{
		cacheKye = cacheKye + vid
	}
	beego.Info("神镜头获取视频详情接口["+vid+"]")
	var videoInfo models.VideoInfo
	cacheHandler,errMsg := models.GetCacheHandler()
	if errMsg != nil {
		videoInfo = models.GetByVid(vid)
		beego.Info("数据从表读取：")
		beego.Info(videoInfo)
		//判断结构vid是否为空，不空，设置缓存
		if videoInfo.Vid != 0 {
			models.SetDataIntoCache(cacheKye,videoInfo,60*3)
		}
	}else{
		if cacheHandler.IsExist(cacheKye) {
			fromCacheByte := cacheHandler.Get(cacheKye).([]byte)
			unmarshalErr := json.Unmarshal(fromCacheByte,&videoInfo)
			if unmarshalErr != nil {
				beego.Info("解析有问题")
				beego.Info(nil)
				videoInfo = models.GetByVid(vid)
				beego.Info("数据从表读取：")
				beego.Info(videoInfo)
				//判断结构vid是否为空，不空，设置缓存
				if videoInfo.Vid != 0 {
					models.SetDataIntoCache(cacheKye,videoInfo,60*3)
				}
			}else{
				beego.Info("数据从缓存读取：")
				beego.Info(videoInfo)
			}
		}else {
			videoInfo = models.GetByVid(vid)
			beego.Info("数据从表读取：")
			beego.Info(videoInfo)
			//判断结构vid是否为空，不空，设置缓存
			if videoInfo.Vid != 0 {
				models.SetDataIntoCache(cacheKye,videoInfo,60*3)
			}
		}
	}
	this.Data["json"] = videoInfo
	this.ServeJSON()
}
