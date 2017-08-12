package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"v-api/models"
	"time"
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
	type user struct{
		Name string
		Age int
	}

	userItem := user{
		"test",
		10,
	}

	testItem := user{}

	//models.SetDataIntoCache("data",video)
	cacheHandler,err := models.GetCacheHandler()
	if err == nil {
		jsonData,_ := json.Marshal(userItem)
		cacheHandler.Put("data",jsonData,10 * time.Second)
		fromCacheByte := cacheHandler.Get("data").([]byte)
		unmarshalErr := json.Unmarshal(fromCacheByte,&testItem)
		if unmarshalErr != nil {
			beego.Info(testItem)
		}else {
			beego.Info(testItem)
		}

	}
	vid := this.Input().Get("vid")
	vidInt, _ := strconv.Atoi(vid)

	testMd5String := "videoInfo_getVideoPlayNum1" + vid
	returnVal := models.Md5(testMd5String)
	fmt.Println(returnVal)
	if cacheHandler.IsExist(returnVal) {
		fromCacheByte := cacheHandler.Get(returnVal)
		fmt.Println(fromCacheByte)
	}




	videoDefinitions, _ := models.GetVideoDefinitions(int64(vidInt), false, "1000,1300,350,yuanhua")
	this.Data["json"] = videoDefinitions
	this.ServeJSON()
	/*client := &http.Client{}
	url := fmt.Sprintf(beego.AppConfig.String("videoTranscodeUrl"),vidInt,"1000,1300,350,yuanhua")
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
		fmt.Println("err:",err)
	}
	req.Header.Set("Host", beego.AppConfig.String("videoTranscodeHost"))
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("result err:",err)
	}

	fmt.Println(string(body))
	*/

	/*models.GetTagByVid(int64(vidInt))


	videoInfo := models.GetRawVideo(vid)
	fmt.Println(videoInfo)
	this.Data["json"] = videoInfo
	this.ServeJSON()*/
	/*var vidStrList []string
	vidStrList = models.GetVidsByUid("50000336")
	fmt.Println(vidStrList)*/
}

// @Title ShenJTLive
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /shenjtlive [get]
func (this *VideoController) ShenJTLive() {
	/*versionCode := this.Input().Get("ver")

	yyuid := this.Input().Get("yyuid")
	limit := this.Input().Get("limit")
	page := this.Input().Get("page")*/

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
		if videoInfo.Vid != "" {
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
				if videoInfo.Vid != "" {
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
			if videoInfo.Vid != "" {
				models.SetDataIntoCache(cacheKye,videoInfo,60*3)
			}
		}
	}
	this.Data["json"] = videoInfo
	this.ServeJSON()
}
