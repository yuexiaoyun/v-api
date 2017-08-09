package models

import (
	"github.com/adam-hanna/arrayOperations"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

}

func GetRawUser(uid string) []orm.Params {
	var rawUser []orm.Params
	o := orm.NewOrm()
	// num, err := o.Raw("SELECT u.vid, u.yyuid, v.user_id, u.video_title, u.video_name, u.source_name, u.channel, u.upload_start_time, u.duration, u.cover, v.video_play_sum, v.video_support FROM  upload_list u LEFT JOIN v_video v ON u.vid = v.vid WHERE u.vid=? AND u.status != -9 AND (u.can_play=1 or u.can_play=4)  LIMIT 1", vid).ValuesList(&lists)
	sql := `SELECT u.user_id, u.nickname, u.approve_name, u.sex, u.udb, u.role,
           ex.user_avatar, ex.user_intro, ex.edit_intro, ex.area, ex.user_level,
           ex.user_game_type, ex.user_video_sum, ex.user_play_sum, ex.user_subscribed_sum
           FROM v_user u LEFT OUTER JOIN v_user_extends ex ON u.user_id = ex.user_id
           WHERE u.user_id = ? `
	o.Raw(sql, uid).Values(&rawUser)
	return rawUser
}
