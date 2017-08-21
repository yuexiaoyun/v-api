package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
	"strings"
)

func init() {

}

type VideoDefinition struct {
	Size       string `json:"size"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	Format     string `json:"format"`
	Duration   string `json:"duration"`
	Definition string `json:"definition"`
	Url        string `json:"url"`
	M3u8       string `json:"m3u8"`
}

type RawVideoInfo struct {
	Vid             int64 `json:"vid"`
	Yyuid           int64 `json:"yyuid"`
	UserId          int64 `json:"user_id"`
	VideoTitle      string
	VideoName       string
	SourceName      string
	Channel         string
	UploadStartTime int64
	Duration        string
	Cover           string
	VideoPlayNum    int64 `orm:"column(video_play_sum)"`
	VideoSupport    int64 `orm:"column(video_support)"`
}

func GetRawVideo(vid string) RawVideoInfo {
	var rawVideo RawVideoInfo
	o := orm.NewOrm()
	sql := `SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid=? AND u.status != -9 AND (u.can_play=1 or u.can_play=4)  LIMIT 1`
	o.Raw(sql, vid).QueryRow(&rawVideo)
	return rawVideo
}

func GetRawVideoByList(vidList []int) []RawVideoInfo {
	var rawVideo []RawVideoInfo
	beego.Info(vidList)
	vidListLen := len(vidList)
	beego.Info(vidListLen)
	var vidListStrSlice []string
	for i := 0; i < vidListLen; i++ {
		vidListStrSlice = append(vidListStrSlice, strconv.Itoa(vidList[i]))
	}
	beego.Info(vidListStrSlice)
	vidListStr := strings.Join(vidListStrSlice, ",")
	beego.Info(vidListStr)
	o := orm.NewOrm()
	sql := `SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid in (?) AND u.status != -9 AND (u.can_play=1 or u.can_play=4)`
	num, err := o.Raw(sql, vidListStr).QueryRows(&rawVideo)
	if err == nil {
		beego.Info("rawVideo nums: ", num)
	}
	return rawVideo
}

/**
	*返回数组[vid,vid,vid,...]格式
	*获取自己上传的视频
**/
func GetVidsByUid(uid string, limit int, page int) []int {
	var vidList []orm.Params
	o := orm.NewOrm()
	start := (page - 1) * limit
	sql := `SELECT u.vid FROM upload_list u LEFT JOIN v_video v ON u.vid=v.vid WHERE u.status!=-9 AND u.can_play=1 AND v.user_id=? limit ?,?`
	o.Raw(sql, uid, start, limit).Values(&vidList)
	var vidIntList []int
	for _, vidMap := range vidList {
		vid := fmt.Sprint(vidMap["vid"])
		vidInt, err := strconv.Atoi(vid)
		if err == nil {
			vidIntList = append(vidIntList, vidInt)
		} else {
			beego.Error(err)
			beego.Error("类型判断失败")
		}

	}
	return vidIntList
}

/**
	*返回数组[vid,vid,vid,...]格式
	*获取账号下的打点视频
**/
func GetDotVidByUid(uid string, limit int, page int) []int {
	var vidList []orm.Params
	o := orm.NewOrm()
	start := (page - 1) * limit
	sql := `SELECT vid FROM upload_list WHERE  yyuid =? and source_client in (14,16) and can_play=1 and status!=-9  ORDER BY upload_start_time DESC LIMIT ?,?`
	o.Raw(sql, uid, start, limit).Values(&vidList)

	var vidIntList []int
	for _, vidMap := range vidList {
		vid := fmt.Sprint(vidMap["vid"])
		vidInt, err := strconv.Atoi(vid)
		if err == nil {
			vidIntList = append(vidIntList, vidInt)
		} else {
			beego.Error(err)
			beego.Error("类型判断失败")
		}
	}
	return vidIntList
}

/*func GetDiyVidByUid(uid string, limit int, page int) []string {
	var vidList []orm.Params
	o := orm.NewOrm()
	start := (page - 1) * limit
	sql := `SELECT vid FROM dot_user_video  WHERE yyuid=? ORDER BY ctime DESC LIMIT ?,?`
	o.Raw(sql, uid, start, limit).Values(&vidList)
	var vidStrList []string
	for _, vidMap := range vidList {
		if vid, ok := vidMap["vid"].(string); ok == true {
			vidStrList = append(vidStrList, vid)
		}
	}
	return vidStrList
}*/

//获取转码信息：视频播放地址，分别率，宽高，大小等信息
func GetVideoDefinitions(vid int64, needAll bool, order string) ([]VideoDefinition, string) {

	cacheKey := VIDEODEFINITIONS
	cacheKey = cacheKey + strconv.Itoa(int(vid))
	cacheHandler, errMsg := GetCacheHandler()
	var videoDefinitions []VideoDefinition
	var status string
	if errMsg == nil {
		if _, _, e := cacheHandler.Get(cacheKey, &videoDefinitions); e != nil {
			videoDefinitions, status = GetVideoDefinitionsFromHost(vid, needAll, order)
			beego.Info("[GetVideoDefinitions]数据从表读取：")
			beego.Info(videoDefinitions)
			//判断结构vid是否为空，不空，设置缓存
			if len(videoDefinitions) != 0 {
				SetDataIntoCache(cacheHandler, cacheKey, videoDefinitions, VIDEODEFINITIONS_TIMEOUT)
			}
		} else {
			beego.Info("[GetVideoDefinitions]数据从缓存读取：")
			beego.Info(videoDefinitions)
		}

	} else {
		beego.Error("[GetVideoByUid]获取缓存句柄失败")
		videoDefinitions, status = GetVideoDefinitionsFromHost(vid, needAll, order)
	}

	return videoDefinitions, status
}

func GetVideoDefinitionsFromHost(vid int64, needAll bool, order string) ([]VideoDefinition, string) {
	url := fmt.Sprintf(beego.AppConfig.String("videoTranscodeUrl"), vid, order)
	req := httplib.Get(url)
	req.Debug(true)
	ret, err := req.String()
	if err != nil {
		fmt.Println(err)
	}
	code := gjson.Get(ret, "code")
	data := gjson.Get(ret, "result")
	var videoDefinitions []VideoDefinition
	if code.String() == "1" {
		data.ForEach(func(key, value gjson.Result) bool {
			var videoDefinition VideoDefinition
			gjson.Unmarshal([]byte(value.String()), &videoDefinition)
			reg := regexp.MustCompile(`(http://huya-w)(.*)(.huya.com)(.*)`)
			result := fmt.Sprint(reg.ReplaceAllString(videoDefinition.Url, "${1}10${3}${4}"))
			videoDefinition.Url = result
			videoDefinitions = append(videoDefinitions, videoDefinition)
			return true // keep iterating
		})
		if len(videoDefinitions) != 0 {
			if !needAll {
				var videoDefinition VideoDefinition
				videoDefinition = videoDefinitions[0]
				for _, videoDefinitionItem := range videoDefinitions {
					if videoDefinitionItem.Format == "m3u8" && videoDefinitionItem.Definition == videoDefinition.Definition {
						videoDefinition.M3u8 = videoDefinitionItem.Url
					}
				}
				return []VideoDefinition{videoDefinition}, "ok"
			}
			return videoDefinitions, "ok"
		} else {
			return videoDefinitions, "no data"
		}
	} else {
		return videoDefinitions, "no data"
	}
}

func GetVideoCategory(channel string) string {
	switch channel {
	case "vhuyalol":
		return "英雄联盟"
	case "vhuyawzry":
		return "王者荣耀"
	case "vhuyaball":
		return "球球大作战"
	case "vhuyacfm":
		return "CFM"
	case "vhuyamc":
		return "我的世界"
	case "vhuyadnf":
		return "地下城与勇士"
	case "vhuyablizzard":
		return "暴雪游戏"
	case "vhuyayule":
		return "娱乐"
	case "vhuyapc":
		return "单机游戏"
	default:
		return ""
	}
}
