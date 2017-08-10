package models

import (
	"github.com/adam-hanna/arrayOperations"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	//"strings"
)

func init() {

}


type VideoDefinition struct{
	Size string `json:"size"`
	Width string `json:"width"`
	Height string `json:"height"`
	Format string `json:"format"`
	Duration string `json:"duration"`
	Definition string `json:"definition"`
	Url string `json:"url"`
	M3u8 string `json:"m3u8"`
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
	VideoPlayNum    int64
	VideoSupport    int64
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

func GetVideoDefinitions(vid int64,needAll bool,order string) /*[]VideoDefinition*/ {
	//TODO 缓存
	/*client := &http.Client{}
	url := fmt.Sprintf(beego.AppConfig.String("videoTranscodeUrl"),vid,order)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle error
	}
	req.Header.Set("Host", beego.AppConfig.String("videoTranscodeHost"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))*/
}
