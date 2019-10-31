package controllers

import (
	"birthday_server/models/stat"
	"github.com/astaxie/beego"
	"log"
)

type GuestStatController struct {
	beego.Controller
}

func (gss *GuestStatController) Show() {
	log.Printf("server for GuestStatController: req=%+v", gss.Input())
	//version 1
	//gss.Data["json"] = stat.GuestMoneyContribution()
	//gss.ServeJSON()
	gss.Data["Resp"] = stat.GuestMoneyContribution()
	gss.TplName = "stat.html"
}
