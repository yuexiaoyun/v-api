package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"v-api/models"
	_ "v-api/routers"
	_ "github.com/astaxie/beego/cache/memcache"

)

func init() {
	models.RegisterDB()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		/*beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"*/
		orm.Debug = true
	}
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

	beego.Run()
}
