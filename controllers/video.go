package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"v-api/models"
	"time"
	"encoding/json"
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
	//versionCode := this.Input().Get("ver")

	/*vid := this.Input().Get("vid")

	models.GetByVid(vid)*/

	vid := this.Input().Get("vid")
	videoInfo := models.GetByVid(vid)

	models.SetDataIntoCache("data",videoInfo,10)
	cacheHandler,_ := models.GetCacheHandler()



	fromCacheByte := cacheHandler.Get("data").([]byte)
	var testItem models.VideoInfo
	unmarshalErr := json.Unmarshal(fromCacheByte,&testItem)
	if unmarshalErr != nil {
		beego.Info(nil)
	}else {
		beego.Info(testItem)
	}


	this.Data["json"] = videoInfo
	this.ServeJSON()
	/*channel := this.Input().Get("channel")
	fmt.Println("channel:",channel)
	game,status := models.GetGameInfoByChannel(channel)
	fmt.Println("status:",status)
	fmt.Println("Game:",game)*/
}
