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
	"strconv"
)

func init() {

}

type UserInfo struct {
	user_id           string
	user_avatar       string
	user_nickname     string
	user_homepage     string
	user_channel      string
	user_channel_name string
	user_video_sum    int
	user_play_sum     int
	user_desc         string
}

func GetRawUser(uid int64) (UserInfo, string) {
	cacheKey := USERINFO
	cacheKey = cacheKey + strconv.Itoa(int(uid))
	cacheHandler, errMsg := GetCacheHandler()
	var userInfo UserInfo
	var status string
	if errMsg != nil {
		userInfo,status = GetRawUserFromDB(uid)
		beego.Info("数据从表读取：")
		beego.Info(userInfo)
		//判断结构vid是否为空，不空，设置缓存
		if userInfo.user_id != "" {
			SetDataIntoCache(cacheHandler,cacheKey, userInfo, USERINFO_TIMEOUT)
		}
	} else {
		if _, _, e := cacheHandler.Get(cacheKey, &userInfo); e != nil {
			userInfo,status = GetRawUserFromDB(uid)
			beego.Info("数据从表读取：")
			beego.Info(userInfo)
			//判断结构vid是否为空，不空，设置缓存
			if userInfo.user_id != ""{
				SetDataIntoCache(cacheHandler,cacheKey, userInfo, USERINFO_TIMEOUT)
			}
		}else{
			status = "ok"
		}
	}
	return userInfo,status
}
func GetRawUserFromDB(uid int64) (UserInfo, string) {
	var rawUser []orm.Params
	o := orm.NewOrm()
	sql := `SELECT u.user_id, u.nickname, u.approve_name, u.sex, u.udb, u.role,
           ex.user_avatar, ex.user_intro, ex.edit_intro, ex.area, ex.user_level,
           ex.user_game_type, ex.user_video_sum, ex.user_play_sum, ex.user_subscribed_sum
           FROM v_user u LEFT OUTER JOIN v_user_extends ex ON u.user_id = ex.user_id
           WHERE u.user_id = ? limit 1`
	num, err := o.Raw(sql, uid).Values(&rawUser)
	var retRawUserInfo UserInfo
	if err == nil && num > 0 {
		for _, userInfo := range rawUser {
			var userChannel string
			userGameType := userInfo["user_game_type"]
			if str, ok := userGameType.(string); ok {
				if strings.Index(str,",") >= 0{
					userChannel = strings.Split(str, ",")[1]
				}else{
					userChannel = ""
				}
			} else {
				userChannel = ""
			}
			gameInfo, status := GetGameInfoByChannel(userChannel)
			var userChannelName string
			if status == "ok" {
				userChannelName = gameInfo.FullName
			}
			userVideoSum, _ := userInfo["user_video_sum"].(int)
			userPlaySum, _ := userInfo["user_play_sum"].(int)
			userAvatar := "http://v.huya.com/style/img/editor-avatar.gif"
			if userInfo["user_avatar"] != nil && userInfo["user_avatar"] != ""{
				userAvatar = fmt.Sprint(userInfo["user_avatar"])
			}
			userId := fmt.Sprint(userInfo["user_id"])
			if  userId == ""{
				userId = strconv.Itoa(int(uid))
			}
			userNickname := fmt.Sprint(userInfo["nickname"])
			if userNickname == ""{
				userNickname = "虎牙小编"
			}

			retRawUserInfo = UserInfo{
				user_id:userId,
				user_avatar:userAvatar,
				user_nickname:userNickname,
				user_homepage:beego.AppConfig.String("baseUrl") + "/u/" + fmt.Sprint(userInfo["user_id"]),
				user_channel: userChannel,
				user_channel_name:userChannelName,
				user_video_sum:userVideoSum,
				user_play_sum:userPlaySum,
				user_desc:fmt.Sprint(userInfo["edit_intro"]),
			}
			break
		}
		return retRawUserInfo, "ok"
	}else{
		retRawUserInfo = UserInfo{
			user_id:strconv.Itoa(int(uid)),
			user_avatar:"http://v.huya.com/style/img/editor-avatar.gif",
			user_nickname:"虎牙小编",
			user_homepage:beego.AppConfig.String("baseUrl") + "/u/" + strconv.Itoa(int(uid)),
			user_channel: "",
			user_channel_name:"",
			user_video_sum:0,
			user_play_sum:0,
			user_desc:"",
		}
		return retRawUserInfo, "ok"
	}

}
