package models

import (
	//"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/astaxie/beego"
	//"strings"
	"github.com/astaxie/beego/orm"
	//"strings"
	//"github.com/astaxie/beego"
	"strings"
	"github.com/astaxie/beego"
	"fmt"
)

func init() {

}

type UserInfo struct{
	user_id string
	user_avatar string
	user_nickname string
	user_homepage string
	user_channel string
	user_channel_name string
	user_video_sum int
	user_play_sum int
	user_desc string
}

func GetRawUser(uid int64) (UserInfo,string) {
	var rawUser []orm.Params
	o := orm.NewOrm()
	// num, err := o.Raw("SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid=? AND u.status != -9 AND (u.can_play=1 or u.can_play=4)  LIMIT 1", vid).ValuesList(&lists)
	sql := `SELECT u.user_id, u.nickname, u.approve_name, u.sex, u.udb, u.role,
           ex.user_avatar, ex.user_intro, ex.edit_intro, ex.area, ex.user_level,
           ex.user_game_type, ex.user_video_sum, ex.user_play_sum, ex.user_subscribed_sum
           FROM v_user u LEFT OUTER JOIN v_user_extends ex ON u.user_id = ex.user_id
           WHERE u.user_id = ? limit 1`
	num,err := o.Raw(sql, uid).Values(&rawUser)
	var retRawUserInfo UserInfo
	if err == nil && num>0 {
		for _,userInfo := range rawUser{
			var userChannel string
			userGameType := userInfo["user_game_type"]
			if str, ok := userGameType.(string); ok {
				userChannel = strings.Split(str,",")[1]
			} else {
				userChannel = ""
			}
			gameInfo,status := GetGameInfoByChannel(userChannel)
			var userChannelName string
			if status == "ok" {
				userChannelName = gameInfo.FullName
			}
			userVideoSum,_ := userInfo["user_video_sum"].(int)
			userPlaySum,_ := userInfo["user_play_sum"].(int)
			userAvatar := "http://v.huya.com/style/img/editor-avatar.gif"
			if  userInfo["user_avatar"]!=nil{
				userAvatar = fmt.Sprint(userInfo["user_avatar"])
			}
			retRawUserInfo = UserInfo{
				user_id:fmt.Sprint(userInfo["user_id"]),
				user_avatar:userAvatar,
				user_nickname:fmt.Sprint(userInfo["nickname"]),
				user_homepage:beego.AppConfig.String("baseUrl")+"/u"+fmt.Sprint(userInfo["user_id"]),
				user_channel: userChannel,
				user_channel_name:userChannelName,
				user_video_sum:userVideoSum,
				user_play_sum:userPlaySum,
				user_desc:fmt.Sprint(userInfo["edit_intro"]),
			}
			break
			/*test := strings.Split(string(userGameType),",")
			fmt.Println(test)*/
			/*
			var userChannel string
			if assert {
				userChannel = strings.Split(userInfo["user_game_type"],",")[0]
			}

			gameInfo,status := GetGameInfoByChannel(userChannel)
			var userChannelName string
			if status == "ok" {
				userChannelName = gameInfo.FullName
			}
			retRawUserInfo = UserInfo{
				user_id:userInfo["user_id"],
				user_avatar:userInfo["user_avatar"],
				user_nickname:userInfo["nickname"],
				user_homepage:beego.AppConfig.String("baseUrl")+"/u"+userInfo["user_id"],
				user_channel: userChannel,
				user_channel_name:userChannelName,
				user_video_sum:userInfo["user_video_sum"],
				user_play_sum:userInfo["user_play_sum"],
				user_desc:userInfo["edit_intro"],
			}
			break*/


		}
		return retRawUserInfo,"ok"
	}
	return UserInfo{},"err:no data"
}
