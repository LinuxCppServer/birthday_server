package log

import (
	"birthday_server/models/json_def"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

const logPath = "bin/log/bds_log.log"

var Glog *logs.BeeLogger

func SysLogSetup() {
	//fileName := json_def.GServer.LogPath + "bs_" + json_def.TimeStampToString(time.Now().Unix(), 0) + ".log"
	fileName := json_def.GServer.LogPath + "bs.log"
	mylog := logs.NewLogger(10000)
	jsonConfigMap := make(map[string]interface{})
	jsonConfigMap["filename"] = fileName
	jsonConfigMap["maxlines"] = 10000
	jsonConfigMap["maxsize"] = 10240
	if jsonConfig, err := json.Marshal(jsonConfigMap); err != nil {
		panic(err)
	} else {
		mylog.SetLogger("file", string(jsonConfig)) // 设置日志记录方式：本地文件记录
		mylog.SetLevel(logs.LevelDebug)             // 设置日志写入缓冲区的等级
		mylog.EnableFuncCallDepth(true)             // 输出log时能显示输出文件名和行号（非必须）
		Glog = mylog
		Glog.Debug("log set up now...")
	}
}
