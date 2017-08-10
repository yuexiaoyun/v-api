package models

import (
	// "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	// "time"s
	"strconv"
	"github.com/astaxie/beego"
	"time"
	"fmt"
	"strings"
)

func init() {

}


type VideoInfo struct {
	UserId       string `json:"user_id"`
	UserAvatar   string `json:"user_avatar"`
	UserNickname string `json:"user_nickname"`
	UserHomepage string `json:"user_homepage"`
	Vid               string `json:"vid"`
	VideoTitle       string `json:"video_title"`
	VideoSubtitle    string `json:"video_subtitle"`
	VideoCover       string `json:"video_cover"`
	VideoCover375x375       string `json:"video_cover_375_375"`
	VideoBigCover      string `json:"video_big_cover"`
	VideoPlayNum    int64 `json:"video_play_num"`
	VideoCommentNum int64 `json:"video_comment_num"`
	VideoDuration    string `json:"video_duration"`
	VideoUrl         string `json:"video_url"`
	VideoUploadTime string `json:"video_upload_time"`
	VideoChannel     string `json:"video_channel"`
	VideoTags        string `json:"video_tags"`
	VideoDefinitions []VideoDefinition `json:"video_definitions"`
}



type RetUserInfo struct {
	user_id       string
	user_avatar   string
	user_nickname string
	user_homepage string
}

func getTitle(rawVideo RawVideoInfo) string {
	var title string
	title = rawVideo.VideoTitle
	return title
}

func getVideoChannel(rawVideo RawVideoInfo) string{
	channel := rawVideo.Channel
	if channel == ""{
		channel = "vhuyaunknown"
	}
	return channel
}

//TODO 年周的计算方法有点繁琐，找时间改进
func getVideoCover(rawVideo RawVideoInfo) string {
	var cover string
	if len(rawVideo.Cover) == 0 {
		yearStr,weekStr := getYearAndWeek(rawVideo.UploadStartTime)
		cover = beego.AppConfig.String("videoCoverPotocal")+beego.AppConfig.String("videoCoverDomain")+"/"+yearStr+weekStr+"/"+strconv.FormatInt(rawVideo.Vid,10)+"/4-220x124.jpg"
	} else {
		cover = rawVideo.Cover
	}
	return cover
}


//TODO 年周的计算方法有点繁琐，找时间改进
func getVideoCoverBySize(rawVideo RawVideoInfo,width string, height string) string {
	var cover string
	yearStr,weekStr := getYearAndWeek(rawVideo.UploadStartTime)
	cover = beego.AppConfig.String("videoCoverPotocal")+beego.AppConfig.String("videoCoverDomain")+"/"+yearStr+weekStr+"/"+strconv.FormatInt(rawVideo.Vid,10)+"/4-"+width+"x"+height+".jpg"
	return cover
}

//TODO 年周的计算方法有点繁琐，找时间改进
func getVideoBigCover(rawVideo RawVideoInfo) string {
	var cover string
	yearStr,weekStr := getYearAndWeek(rawVideo.UploadStartTime)
	cover = beego.AppConfig.String("videoCoverPotocal")+beego.AppConfig.String("videoCoverDomain")+"/"+yearStr+weekStr+"/"+strconv.FormatInt(rawVideo.Vid,10)+"/4-640x360.jpg"
	return cover
}

func getYearAndWeek(timestamp int64) (string,string){
	format := "2006-01-02 15:04:05"
	t, _ := time.Parse(format, time.Unix(timestamp,0).Format(format))
	year, week := t.ISOWeek()
	year = year - 2000
	var yearStr string
	var weekStr string
	if year < 10{
		yearStr = "0" + strconv.Itoa(year)
	}else{
		yearStr = strconv.Itoa(year)
	}
	if week < 10{
		weekStr = "0" + strconv.Itoa(week)
	}else{
		weekStr = strconv.Itoa(week)
	}

	return yearStr,weekStr
}


func getVideoInfo(rawVideo RawVideoInfo) VideoInfo{
	rawUser,status := GetRawUser(rawVideo.Yyuid)
	var retUserInfo RetUserInfo = RetUserInfo{}
	if status == "ok" {
		retUserInfo = RetUserInfo{
			user_id:rawUser.user_id,
			user_avatar:rawUser.user_avatar,
			user_nickname:rawUser.user_nickname,
			user_homepage:rawUser.user_homepage,
		}
	}

	videoInfo := VideoInfo{
		Vid:            strconv.FormatInt(rawVideo.Vid,10),
		VideoTitle:    getTitle(rawVideo),
		VideoSubtitle: getTitle(rawVideo),
		VideoCover:    getVideoCover(rawVideo),
		VideoPlayNum:rawVideo.VideoPlayNum,
		VideoCommentNum:rawVideo.VideoSupport,
		VideoDuration:rawVideo.Duration,
		VideoBigCover:getVideoBigCover(rawVideo),
		VideoCover375x375:getVideoCoverBySize(rawVideo,"375","375"),
		VideoUrl:getVideoUrl(rawVideo),
		VideoUploadTime:getVideoUploadTime(rawVideo),
		VideoChannel:getVideoChannel(rawVideo),
		VideoTags:getVideoTags(rawVideo),
		UserId:retUserInfo.user_id,
		UserAvatar:retUserInfo.user_avatar,
		UserNickname:retUserInfo.user_nickname,
		UserHomepage:retUserInfo.user_homepage,
	}
	return videoInfo
}

func getVideoUrl(rawVideo RawVideoInfo) string{
	return beego.AppConfig.String("baseUrl") +"/play/" + strconv.FormatInt(rawVideo.Vid,10) + ".html"
}

func getVideoUploadTime(rawVideo RawVideoInfo) string {
	return fmt.Sprint(time.Unix(rawVideo.UploadStartTime,0).Format("2006-01-02 03:04:05"))
}
func getVideoTags(rawVideo RawVideoInfo) string{
	videoTags := GetTagByVid(rawVideo.Vid)
	var tags []string
	for _,tag := range videoTags{
		tags = append(tags,tag.Tag)
	}
	if len(tags) == 0 {
		return ""
	}else{
		return strings.Join(tags,",")
	}
}

func getUserInfo(uid int64)  RetUserInfo{
	rawUser,status := GetRawUser(uid)
	var retUserInfo RetUserInfo = RetUserInfo{}
	if status == "ok" {
		retUserInfo = RetUserInfo{
			user_id:rawUser.user_id,
			user_avatar:rawUser.user_avatar,
			user_nickname:rawUser.user_nickname,
			user_homepage:rawUser.user_homepage,
		}
		return retUserInfo
	}
	return retUserInfo
}

func GetByVid(vid string) VideoInfo {
	//TODO 缓存
	rawVideo := GetRawVideo(vid)
	videoInfo:= getVideoInfo(rawVideo)

	fmt.Println(videoInfo)

	/*rawUser := GetRawUser(int(video.Yyuid))
	fmt.Println(rawUser)*/
	return videoInfo
}