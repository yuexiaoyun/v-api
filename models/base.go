package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/cache"
	"encoding/json"
	"time"
)

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

func GetCacheHandler() (adapter cache.Cache, err error){
	return cache.NewCache("memcache", `{"conn":"`+beego.AppConfig.String("memcache_host_1")+`:`+beego.AppConfig.String("memcache_port_1")+`;`+beego.AppConfig.String("memcache_host_2")+`:`+beego.AppConfig.String("memcache_port_2")+`"}`)
}

func SetDataIntoCache(key string,data interface{},timeout int64){
	cacheHandler,err := GetCacheHandler()
	if err == nil {
		jsonData,_ := json.Marshal(data)
		cacheHandler.Put(key,jsonData,time.Duration(timeout) * time.Second)
	}else {
		beego.Error("缓存设置数据出错")
		beego.Error(err)
	}

}

