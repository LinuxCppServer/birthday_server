//此代码是宾客管理模块相关的后端接口
package controllers

import (
	"birthday_server/models/json_def"
	. "birthday_server/models/log"
	"birthday_server/models/mysql"
	. "birthday_server/models/notify"
	"birthday_server/models/stat"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"sort"
	"strconv"
	"time"
)

type GuestManageController struct {
	beego.Controller
}

//add
func (gmc *GuestManageController) AddNewGuest() {
	Glog.Debug("serve for AddNewGuest: req=%+v", gmc.Input())
	if len(gmc.Input()) < 2 {
		Glog.Error("AddNewGuest failed: parameters aren't enough")
		gmc.Data["json"] = json_def.PackJsonAck(1, "增加宾客记录参数不少于2个", nil)
		gmc.ServeJSON()
		return
	}
	var money, attend_count int
	var name string
	var err error
	name = gmc.GetString("name")
	if money, err = gmc.GetInt("money"); err != nil {
		Glog.Error("money isn't int-like: err=%s", err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(2, "增加宾客记录礼金参数应是整型数", nil)
		gmc.ServeJSON()
		return
	}
	if attend_count, err = gmc.GetInt("attend_count"); err != nil {
		Glog.Error("attend count isn't int-like: err=%s", err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(3, "增加宾客记录出席人数应是整型数", nil)
		gmc.ServeJSON()
		return
	}
	entry_time := json_def.TimeStampToString(time.Now().Unix(), 0)
	Glog.Error("准备增加宾客记录: name=%s money=%d attend_count=%d entry_time=%s",
		name, money, attend_count, entry_time)
	//versin1: sql := "insert into t_guest_money(name, money, attend_count) values(?,?,?)"
	sql := "insert into t_guest_money(name, money, attend_count, entry_time) values(?,?,?,?)"
	if _, err = mysql.Cmd(mysql.MYSQL_CMD_EXEC, sql, name, money, attend_count, entry_time); err != nil {
		Glog.Error("insert into mysql for add new guest failed: err=%s", err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(4, fmt.Sprintf("增加宾客礼金记录失败: %s", err.Error()), nil)
		gmc.ServeJSON()
		return
	}
	msg := fmt.Sprintf("恭喜;><: 有一笔新礼金入袋!\n姓名：%s\n礼金：%d\n出席人数：%d", name, money, attend_count)
	GMsgNotify.NotifyEveryOne(msg)
	if info := stat.SubTotalInfo(); info != "" {
		GStatNotify.NotifyEveryOne(info)
	}
	gmc.Data["json"] = json_def.PackJsonAck(0, "ok", fmt.Sprintf("增加宾客记录成功：\n姓名：%s\n礼金：%d\n出席人数：%d", name, money, attend_count))
	gmc.ServeJSON()
}

//del
func (gmc *GuestManageController) DelGuestById() {
	Glog.Debug("serve for DelGuestById: req=%+v", gmc.Input())
	if len(gmc.Input()) <= 0 {
		Glog.Error("DelGuestById failed: parameters aren't enough")
		gmc.Data["json"] = json_def.PackJsonAck(5, "删除宾客记录青提供id", nil)
		gmc.ServeJSON()
		return
	}

	var id int = -1
	var err error
	id, err = gmc.GetInt("id")
	if err != nil {
		Glog.Error("id isn't int-like")
		gmc.Data["json"] = json_def.PackJsonAck(6, "宾客id应是整型数", nil)
		gmc.ServeJSON()
		return
	}

	querySql := fmt.Sprintf("select * from t_guest_money where id=%d", id)
	if rows, err := mysql.Cmd(mysql.MYSQL_CMD_QUERY, querySql); err != nil {
		Glog.Error("query sql isn't correct: sql=%s", querySql)
		gmc.Data["json"] = json_def.PackJsonAck(7, "internel error", nil)
		gmc.ServeJSON()
		return
	} else {
		delHintInfo, basicInfo := "", "\nid:%d  姓名:%s  礼金:%d  出席人数:%d  录入时间:%s\n"
		var (
			_name, _entry_time         string
			_id, _money, _attend_count int
			rec_count                  int
		)

		for rows.Next() {
			rec_count++
			if e := rows.Scan(&_id, &_name, &_money, &_attend_count, &_entry_time); e != nil {
				Glog.Error("scan record failed: id=%d name=%s monney=%d attend_count=%d entry_time=%s err=%s",
					_id, _name, _money, _attend_count, _entry_time, e.Error())
				continue
			}
			delHintInfo += fmt.Sprintf(basicInfo, _id, _name, _money, _attend_count, _entry_time)
		}
		if rec_count == 0 {
			Glog.Error("the id want to delete isn't exists: id=%d", id)
			gmc.Data["json"] = json_def.PackJsonAck(8, fmt.Sprintf("欲删除的宾客记录不存在: id=%d", id), nil)
			gmc.ServeJSON()
			return
		}
		delHintInfo += "以上宾客记录已被删除!"
		GMsgNotify.NotifyEveryOne(delHintInfo)
	}

	ssql := "delete from t_guest_money where id=?"
	if _, err = mysql.Cmd(mysql.MYSQL_CMD_EXEC, ssql, id); err != nil {
		Glog.Error("delete from mysqwl failed: id=%d err=%s", id, err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(9, fmt.Sprintf("删除记录失败: err=%s", err.Error()), nil)
		gmc.ServeJSON()
		return
	}
	if info := stat.SubTotalInfo(); info != "" {
		GStatNotify.NotifyEveryOne(info)
	}
	gmc.Data["json"] = json_def.PackJsonAck(0, "ok", "删除成功")
	gmc.ServeJSON()
}

//mod
func (gmc *GuestManageController) ModifyGuestRecord() {
	Glog.Debug("ModifyGuestRecord: req=%+v", gmc.Input())
	if len(gmc.Input()) < 2 {
		Glog.Error("ModifyGuestRecord failed: parameters aren't enough")
		gmc.Data["json"] = json_def.PackJsonAck(10, "修改宾客记录至少含有2个参数", nil)
		gmc.ServeJSON()
		return
	}
	var id, money, attend_count int = -1, 0, 0
	var err error
	if id, err = gmc.GetInt("id"); err != nil {
		Glog.Error("id isn't int-like")
		gmc.Data["json"] = json_def.PackJsonAck(11, "修改宾客记录的id参数应是整型数", nil)
		gmc.ServeJSON()
		return
	}

	ssql, name := "update t_guest_money set ", gmc.GetString("name")
	if name != "" {
		ssql += fmt.Sprintf("name=\"%s\",", name)
	}
	if newMoneyStr := gmc.Input().Get("money"); newMoneyStr != "" {
		if money, err = gmc.GetInt("money"); err != nil {
			Glog.Error("money for ModifyGuest must be int-like")
			gmc.Data["json"] = json_def.PackJsonAck(12, "修改宾客的新礼金参数应是整型数", nil)
			gmc.ServeJSON()
			return
		}
		ssql += fmt.Sprintf("money=%d,", money)
	}
	if newAttendCount := gmc.Input().Get("count"); newAttendCount != "" {
		if attend_count, err = gmc.GetInt("count"); err != nil {
			Glog.Error("attend_count for ModifyGuest must be int-like")
			gmc.Data["json"] = json_def.PackJsonAck(13, "修改宾客的新出席人数应是整型数", nil)
			gmc.ServeJSON()
			return
		}
		ssql += fmt.Sprintf("attend_count=%d,", attend_count)
	}
	ssql = ssql[:len(ssql)-1]
	ssql += fmt.Sprintf("  where id=%d", id)
	Glog.Debug("will modify guest record: id=%d sql=%s", id, ssql)

	querySql, modHintInfo := fmt.Sprintf("select * from t_guest_money where id=%d", id), "宾客记录修改详情:"
	if rows, err := mysql.Cmd(mysql.MYSQL_CMD_QUERY, querySql); err != nil {
		Glog.Error("query sql isn;t incorrect: err=%s", err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(14, err.Error(), nil)
		gmc.ServeJSON()
		return
	} else {
		if rows == nil {
			Glog.Error("the quest isn't exists: id=%d", id)
			gmc.Data["json"] = json_def.PackJsonAck(15, fmt.Sprintf("宾客不存在: id=%d", id), nil)
			gmc.ServeJSON()
			return
		}

		var (
			_name, _entry_time         string
			_id, _money, _attend_count int
		)
		for rows.Next() {
			if e := rows.Scan(&_id, &_name, &_money, &_attend_count, &_entry_time); e != nil {
				break
			}
		}
		modHintInfo += fmt.Sprintf("\n修改前: id=%d  姓名: %s  礼金: %d  出席人数: %d",
			_id, _name, _money, _attend_count)
	}

	if _, err = mysql.Cmd(mysql.MYSQL_CMD_EXEC, ssql); err != nil {
		Glog.Error("update mysql failed: err=%s", err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(16, fmt.Sprintf("修改宾客记录失败: err=%s", err.Error()), nil)
		gmc.ServeJSON()
		return
	}

	rows, _ := mysql.Cmd(mysql.MYSQL_CMD_QUERY, querySql)
	var (
		_name, _entry_time         string
		_id, _money, _attend_count int
	)
	for rows.Next() {
		if e := rows.Scan(&_id, &_name, &_money, &_attend_count, &_entry_time); e != nil {
			break
		}
	}
	modHintInfo += fmt.Sprintf("\n修改后: id=%d  姓名: %s  礼金: %d  出席人数: %d",
		_id, _name, _money, _attend_count)
	GMsgNotify.NotifyEveryOne(modHintInfo)
	if info := stat.SubTotalInfo(); info != "" {
		GStatNotify.NotifyEveryOne(info)
	}
	retList := []*json_def.QueryGuestAck{}
	retList = append(retList, &json_def.QueryGuestAck{
		_id,
		_name,
		_money,
		_attend_count,
		_entry_time,
	})
	gmc.Data["json"] = json_def.PackJsonAck(0, "ok", retList)
	gmc.ServeJSON()
}

const (
	query_stub = iota
	query_guest_by_full_name
	query_guest_by_last_name
	query_guest_by_money_range
	query_guest_all
)

func (gmc *GuestManageController) QueryGuests() {
	Glog.Debug("serve for QueryGuests: req=%+v", gmc.Input())
	if len(gmc.Input()) < 2 {
		Glog.Error("QueryGuests failed: parameters aren't enough")
		gmc.Data["json"] = json_def.PackJsonAck(17, "查询宾客记录参数应不少于2个", nil)
		gmc.ServeJSON()
		return
	}
	var option int
	var paraList []string
	var err error
	if option, err = gmc.GetInt("option"); err != nil {
		Glog.Error("option for query isn't int-like")
		gmc.Data["json"] = json_def.PackJsonAck(18, "查询宾客记录的option参数应是整型数", nil)
		gmc.ServeJSON()
		return
	}
	paraList = gmc.GetStrings("paras[]")
	if len(paraList) == 0 {
		Glog.Error("query para list is null")
		gmc.Data["json"] = json_def.PackJsonAck(19, "查询宾客记录参数列表不能为空", nil)
		gmc.ServeJSON()
		return
	}

	var ssql string
	var args []interface{}
	switch option {
	case query_guest_by_full_name:
		ssql = "select * from t_guest_money where name=?"
		Glog.Debug("query guest by full-name: sql=%s", ssql)
		args = append(args, paraList[0])
	case query_guest_by_last_name:
		ssql = fmt.Sprintf(`select * from t_guest_money where name like '%s%%'`, paraList[0])
		Glog.Debug("fuzzly query by last-name: last_name=%s, len(last_name)=%d sql=%s", paraList[0], len(paraList[0]), ssql)
	case query_guest_by_money_range:
		Glog.Debug("fuzzly query by money range: [%s, %s]", paraList[0], paraList[1])
		var min, max int
		if min, err = strconv.Atoi(paraList[0]); err != nil {
			Glog.Error("min range must be int-like")
			gmc.Data["json"] = json_def.PackJsonAck(20, fmt.Sprintf("礼金范围必须是整型数： min=%s", paraList[0]), nil)
			gmc.ServeJSON()
			return
		}
		if max, err = strconv.Atoi(paraList[1]); err != nil {
			Glog.Error("max range must be int-like")
			gmc.Data["json"] = json_def.PackJsonAck(21, fmt.Sprintf("礼金范围必须是整型数： max=%s", paraList[1]), nil)
			gmc.ServeJSON()
			return
		}
		if min < 0 || max < 0 {
			Glog.Error("money range must be >= 0")
			gmc.Data["json"] = json_def.PackJsonAck(22, fmt.Sprintf("礼金范围必须是+整型数： min=%d max=%d", min, max), nil)
			gmc.ServeJSON()
			return
		}
		if min > max {
			min, max = max, min
		}
		ssql = "select * from t_guest_money where money between  ? and ?"
		Glog.Debug("fuzzly query by money range: sql=%s", ssql)
		args = append(args, min, max)
	case query_guest_all:
		Glog.Debug("query all guests, paras->%+v will be ignored", paraList)
		ssql = "select * from t_guest_money"
	case query_stub:
		ssql = fmt.Sprintf("select * from t_guest_money where id=%s", paraList[0])
		Glog.Debug("internel query interface: sql=%s", ssql)
	default:
		Glog.Debug("unkown query option: %d", option)
		gmc.Data["json"] = json_def.PackJsonAck(23, "系统目前暂不支持此种查询", nil)
		gmc.ServeJSON()
		return
	}

	Glog.Debug("ssql=%s\nargs=%+v", ssql, args)
	var rows *sql.Rows
	rows, err = mysql.Cmd(mysql.MYSQL_CMD_QUERY, ssql, args...)
	if err != nil {
		Glog.Error("exec query failed: err=%s", err.Error())
		gmc.Data["json"] = json_def.PackJsonAck(24, "数据库查询失败", nil)
		gmc.ServeJSON()
		return
	}

	list := []*json_def.QueryGuestAck{}
	for rows.Next() {
		name, entry_time := "", ""
		id, money, attend_count := -1, 0, 0
		if err = rows.Scan(&id, &name, &money, &attend_count, &entry_time); err != nil {
			Glog.Error("result scan failed: err=%s", err.Error())
			continue
		}
		Glog.Debug("id=%d name=%s money=%d attend_count=%d entry_time=%s", id, name, money, attend_count, entry_time)
		list = append(list, &json_def.QueryGuestAck{
			id,
			name,
			money,
			attend_count,
			entry_time,
		})
	}
	sort.Sort(json_def.GuestRankHelper(list))
	gmc.Data["json"] = json_def.PackJsonAck(0, "ok", list)
	gmc.ServeJSON()
}
