package main

import (
	. "birthday_server/models/log"
	"birthday_server/models/mysql"
	_ "birthday_server/routers"
	"github.com/astaxie/beego"
)

func main() {
	SysLogSetup()
	mysql.GetDBConn()
	Glog.Debug("server start now...")
	beego.Run()
	mysql.ReleaseDB()
	Glog.Debug("server exited safely...")
}
