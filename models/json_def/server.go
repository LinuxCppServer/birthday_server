package json_def

import (
	"encoding/json"
	"io/ioutil"
)

var GServer struct {
	LogPath              string
	LogLevel             int
	ConfPath             string
	MaxHouseHoldBooked   int
	MaxTableBooked       int
	MaxGuestPersonBooked int
	MoneyPackForBackVal  int
	MoneyPackForTeaVal   int
}

func init() {
	data, err := ioutil.ReadFile("bin/conf/server.json")
	if err != nil {
		panic(err)
	}

	//使用自定义的读取配置文件的方式
	var scMap map[string]interface{}
	if err := json.Unmarshal([]byte(data), &scMap); err == nil {
		GServer.LogLevel = int(scMap["LogLevel"].(float64))
		GServer.LogPath = scMap["LogPath"].(string)
		GServer.ConfPath = "bin/conf"
		GServer.MaxHouseHoldBooked = int(scMap["MaxHouseHoldBooked"].(float64))
		GServer.MaxTableBooked = int(scMap["MaxTableBooked"].(float64))
		GServer.MaxGuestPersonBooked = int(scMap["MaxGuestPersonBooked"].(float64))
		GServer.MoneyPackForBackVal = int(scMap["MoneyPackForBackVal"].(float64))
		GServer.MoneyPackForTeaVal = int(scMap["MoneyPackForTeaVal"].(float64))
	}
}
