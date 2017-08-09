package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/adam-hanna/arrayOperations"
)

func init() {

}

type RawVideoInfo struct {
	Vid             int64
	Yyuid           int64
	UserId          int64
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
	o.Raw("SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid=? AND u.status != -9 AND (u.can_play=1 or u.can_play=4)  LIMIT 1", vid).QueryRow(&rawVideo)
	return rawVideo
}

/**
	*返回数组[vid,vid,vid,...]格式
	*获取自己上传的视频
**/
func GetVidsByUid(uid string) []string {
	var vidList []orm.Params
	o := orm.NewOrm()
	o.Raw("SELECT u.vid FROM upload_list u LEFT JOIN v_video v ON u.vid=v.vid WHERE u.status!=-9 AND u.can_play=1 AND v.user_id=?", uid).Values(&vidList)

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
	o.Raw("SELECT vid FROM upload_list WHERE  yyuid =? and source_client in (14,16) and can_play=1 and status!=-9  ORDER BY upload_start_time DESC LIMIT ?,?", uid, start, limit).Values(&vidList)

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
	o.Raw("SELECT vid FROM dot_user_video  WHERE yyuid=? ORDER BY ctime DESC LIMIT ?,?", uid, start, limit).Values(&vidList)
	var vidStrList []string
	for _, vidMap := range vidList {
		if vid, ok := vidMap["vid"].(string); ok == true {
			vidStrList = append(vidStrList, vid)
		}
	}
	return vidStrList
}


func GetVideoByUid(yyuid string,limit int,page int) {
	/**
	TODO 缓存
	 */
	liveVids := GetDotVidByUid(yyuid,limit,page)
	diyVids := GetDiyVidByUid(yyuid,limit,page)
	uploadVids := GetVidsByUid(yyuid)


	vids,ret := arrayOperations.Union(liveVids,diyVids,uploadVids)
	if ret == "ok"{

	}


}