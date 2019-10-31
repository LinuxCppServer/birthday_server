package routers

import (
	"birthday_server/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/main/debug", &controllers.MainController{}, "*:Debug")
	// root -> login.html
	beego.Router("/", &controllers.MainController{}, "*:Login")
	// /main  -> main.html
	beego.Router("/main", &controllers.MainController{}, "*:Main")
	// /manage -> manage.html
	beego.Router("/manage", &controllers.MainController{}, "*:Manage")
	// /stat -> stat.html
	beego.Router("/stat", &controllers.MainController{}, "*:Stat")

	beego.Router("/manage/add_new_guest", &controllers.GuestManageController{}, "*:AddNewGuest")
	beego.Router("/manage/query_guest", &controllers.GuestManageController{}, "*:QueryGuests")
	beego.Router("/manage/del_guest", &controllers.GuestManageController{}, "*:DelGuestById")
	beego.Router("/manage/mod_guest", &controllers.GuestManageController{}, "*:ModifyGuestRecord")

	beego.Router("/stat/stat_show", &controllers.GuestStatController{}, "*:Show")
}
