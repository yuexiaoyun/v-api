package models

import (
	"github.com/adam-hanna/arrayOperations"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/tidwall/gjson"
	"regexp"
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
	// num, err := o.Raw("SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid=? AND u.status != -9 AND (u.can_play=1 or u.can_play=4)  LIMIT 1", vid).ValuesList(&lists)
	sql := `SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid=? AND u.status != -9 AND (u.can_play=1 or u.can_play=4)  LIMIT 1`
	o.Raw(sql, vid).QueryRow(&rawVideo)
	return rawVideo
}

/**
	*返回数组[vid,vid,vid,...]格式
	*获取自己上传的视频
**/
func GetVidsByUid(uid string) []string {
	var vidList []orm.Params
	o := orm.NewOrm()
	sql := `SELECT u.vid FROM upload_list u LEFT JOIN v_video v ON u.vid=v.vid WHERE u.status!=-9 AND u.can_play=1 AND v.user_id=?`
	o.Raw(sql, uid).Values(&vidList)

	var vidStrList []string
	for _, vidMap := range vidList {
		if vid, ok := vidMap["vid"].(string); ok == true {
			vidStrList = append(vidStrList, vid)
		}
	}
	return vidStrList
}

/**
	*返回数组[vid,vid,vid,...]格式
	*获取账号下的打点视频
**/
func GetDotVidByUid(uid string, limit int, page int) []string {
	var vidList []orm.Params
	o := orm.NewOrm()
	start := (page - 1) * limit
	sql := `SELECT vid FROM upload_list WHERE  yyuid =? and source_client in (14,16) and can_play=1 and status!=-9  ORDER BY upload_start_time DESC LIMIT ?,?`
	o.Raw(sql, uid, start, limit).Values(&vidList)

	var vidStrList []string
	for _, vidMap := range vidList {
		if vid, ok := vidMap["vid"].(string); ok == true {
			vidStrList = append(vidStrList, vid)
		}
	}
	return vidStrList
}

func GetDiyVidByUid(uid string, limit int, page int) []string {
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
}

func GetVideoByUid(yyuid string, limit int, page int) {
	/**
	TODO 缓存
	*/
	liveVids := GetDotVidByUid(yyuid, limit, page)
	diyVids := GetDiyVidByUid(yyuid, limit, page)
	uploadVids := GetVidsByUid(yyuid)

	vids, ret := arrayOperations.Union(liveVids, diyVids, uploadVids)
	if ret {
		fmt.Println(vids)
	}

}
//获取转码信息：视频播放地址，分别率，宽高，大小等信息
func GetVideoDefinitions(vid int64, needAll bool, order string) ([]VideoDefinition, string) {
	//TODO 缓存
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
	/*categoryMap := make(map[string]string,10)
	categoryMap["vhuyalol"] = "英雄联盟"
	categoryMap["vhuyawzry"] = "王者荣耀"
	categoryMap["vhuyaball"] = "球球大作战"
	categoryMap["vhuyacfm"] = "CFM"
	categoryMap["vhuyamc"] = "我的世界"
	categoryMap["vhuyadnf"] = "地下城与勇士"
	categoryMap["vhuyablizzard"] = "暴雪游戏"
	categoryMap["vhuyayule"] = "娱乐"
	categoryMap["vhuyapc"] = "单机游戏"*/
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
