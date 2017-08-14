package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["v-api/controllers:VideoController"] = append(beego.GlobalControllerRouter["v-api/controllers:VideoController"],
		beego.ControllerComments{
			Method: "ShenJTDetail",
			Router: `/shenjtdetail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["v-api/controllers:VideoController"] = append(beego.GlobalControllerRouter["v-api/controllers:VideoController"],
		beego.ControllerComments{
			Method: "ShenJTLive",
			Router: `/shenjtlive`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["v-api/controllers:VideoController"] = append(beego.GlobalControllerRouter["v-api/controllers:VideoController"],
		beego.ControllerComments{
			Method: "TestRouter",
			Router: `/test`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
