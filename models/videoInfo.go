package models

import (
	// "github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// "strconv"
	// "time"
)

func init() {

}

type ShowedInfo struct {
	UserInfo
	VideoInfo
}
type VideoInfo struct {
	vid               string
	video_title       string
	video_subtitle    string
	video_cover       string
	video_play_num    string
	video_comment_num string
	video_duration    string
	video_url         string
	video_upload_time string
	video_channel     string
	video_tags        string
}

type UserInfo struct {
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

func getVideoChannel() {

}

//TODO
func getVideoCover(rawVideo RawVideoInfo) string {
	var cover string
	if len(rawVideo.Cover) == 0 {

	} else {
		cover = rawVideo.Cover
	}
	return cover
}

func getVideoInfo(rawVideo RawVideoInfo) {
	/*videoInfo = VideoInfo{
		vid:            strconv.FormatInt(rawVideo.Vid),
		video_title:    getTitle(rawVideo),
		video_subtitle: getTitle(rawVideo),
		video_cover:    getVideoCover(rawVideo),
		/*video_play_num    string
		video_comment_num:0,
		video_duration    string
		video_url         string
		video_upload_time string
		video_channel     string
		video_tags        string
	}
	*/
}
