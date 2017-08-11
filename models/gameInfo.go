package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
)

func init() {
	orm.RegisterModel(new(GameInfo))
}

type GameInfo struct {
	Id          int `orm:"column(game_id)"`
	ShortName   string `orm:"column(short_name)"`
	FullName    string `orm:"column(full_name)"`
	EnglishName string `orm:"column(english_name)"`
	Channel     string `orm:"column(channel)"`
	GameType    int `orm:"column(game_type)"`
	GameCover   string `orm:"column(game_cover)"`
	GameIntro   string `orm:"column(game_intro)"`
	VideoSum    int `orm:"column(video_sum)"`
	Pinyin      string `orm:"column(pinyin)"`
	IsPublic    int `orm:"column(is_public)"`
}

func (gameInfo *GameInfo) TableName() string {
	return "v_game"
}

func GetGameInfoByChannel(channel string) (GameInfo, string) {
	var gameInfo GameInfo
	gameInfo = GameInfo{Channel:channel}
	o := orm.NewOrm()
	err := o.Read(&gameInfo, "channel")
	if err == orm.ErrNoRows {
		return GameInfo{}, "err:不存在"
	} else if err == orm.ErrMissPK {
		return GameInfo{}, "err:找不到主键"
	} else {
		return gameInfo, "ok"
	}
}
