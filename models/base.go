package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/pangudashu/memcache"
	_ "github.com/go-sql-driver/mysql"
)

const (
	SHENJTLIVE   = "shenjtlive_v1_"
	SHENJTDETAIL = "shenjtdetail_v1_"
	VIDEOINFO    = "videoinfo_vi_"
	USERINFO     = "userinfo_v1_"
	VIDEODEFINITIONS = "video_definitions_v1_"

	SHENJTLIVE_TIMEOUT = 30 * 60
	SHENJTDETAIL_TIMEOUT = 30 * 60
	VIDEOINFO_TIMEOUT = 30 * 60
	VIDEODEFINITIONS_TIMEOUT = 30 * 60
	USERINFO_TIMEOUT = 15 * 60
)


type DatabaseCheck struct {
}

func (dc *DatabaseCheck) Check() error {
	if dc.IsConnected() {
		return nil
	} else {
		return errors.New("can't connect database")
	}
}

func (dc *DatabaseCheck) IsConnected() bool {
	db, err := orm.GetDB("default")
	if err == nil {
		if db.Ping() == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func RegisterDB() {
	dbhost := beego.AppConfig.String("mysqlhost")
	dbport := beego.AppConfig.String("mysqlport")
	dbuser := beego.AppConfig.String("mysqluser")
	dbpassword := beego.AppConfig.String("mysqlpass")
	db := beego.AppConfig.String("mysqldb")
	//注册mysql Driver
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//构造conn连接
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + db + "?charset=utf8"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)
	fmt.Printf("数据库连接成功！%s\n", conn)
}

func GetCacheHandler() (mem *memcache.Memcache, err error) {
	if enableCache, _ := beego.AppConfig.Bool("EnableCache"); enableCache == false {
		err = fmt.Errorf("Cache is disable now")
		return
	}

	s1 := &memcache.Server{Address: beego.AppConfig.String("memcache_host_1")+":"+beego.AppConfig.String("memcache_port_1"), Weight: 50}
	s2 := &memcache.Server{Address: beego.AppConfig.String("memcache_host_2")+":"+beego.AppConfig.String("memcache_port_2"), Weight: 50}
	mem,err = memcache.NewMemcache([]*memcache.Server{s1, s2})
	if err != nil {
		beego.Error("缓存初始化链接失败")
		beego.Error(err)
		return
	}
	return mem,err
	//return cache.NewCache("memcache", `{"conn":"`+beego.AppConfig.String("memcache_host_1")+`:`+beego.AppConfig.String("memcache_port_1")+`;`+beego.AppConfig.String("memcache_host_2")+`:`+beego.AppConfig.String("memcache_port_2")+`"}`)
}

/*func GetCacheHandler() (adapter cache.Cache, err error) {
	if enableCache, _ := beego.AppConfig.Bool("EnableCache"); enableCache == false {
		err = fmt.Errorf("Cache is disable now")
		return
	}

	return cache.NewCache("memcache", `{"conn":"`+beego.AppConfig.String("memcache_host_1")+`:`+beego.AppConfig.String("memcache_port_1")+`;`+beego.AppConfig.String("memcache_host_2")+`:`+beego.AppConfig.String("memcache_port_2")+`"}`)
}*/

func SetDataIntoCache(cacheHandler *memcache.Memcache,key string, data interface{}, timeout uint32) {

	if enableCache, _ := beego.AppConfig.Bool("EnableCache"); enableCache == false {
		beego.Info("不开缓存，无法设置")
		return
	}else{
		cacheHandler.Set(key, data, timeout)
	}
}

func Md5(value string) string {
	data := []byte(value)
	return fmt.Sprintf("%x", md5.Sum(data))
}
