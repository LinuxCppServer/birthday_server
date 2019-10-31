package stat

import (
	"birthday_server/models/json_def"
	. "birthday_server/models/log"
	"birthday_server/models/mysql"
	"database/sql"
	"fmt"
	"sort"
	"time"
)

//todo read from conf
const (
	MAX_HOUSEHOLD_BOOKED    = 5  //80 	//预定户数
	MAX_TABLE_BOOKED        = 16 //预定桌数
	MAX_GUEST_PERSON_BOOKED = 20 //160	//预定人数
	MONEY_PACK_FOR_BACK_VAL = 12 //回礼红包价值
	MONEY_PACK_FOR_TEA_VAL  = 4  //茶红包价值
)

func SubTotalInfo() (info string) {
	var (
		resp             error
		rows             *sql.Rows
		householdsCount  int //户数
		guestPersonCount int //人数
		totalGuestMoney  int //总礼金
		guestAdded       int //赶席总人数
		tableAdded       int //赶席桌数
		moneyTotalAdded  int //需要准备的现金数

		name, entry_time, sql   string = "", " ", "select * from t_guest_money"
		id, money, attend_count int
	)

	if rows, resp = mysql.Cmd(mysql.MYSQL_CMD_QUERY, sql); resp != nil {
		Glog.Error("SubTotalInfo failed: sql=%s err=%s", sql, resp.Error())
		return resp.Error()
	}

	for rows.Next() {
		if err := rows.Scan(&id, &name, &money, &attend_count, &entry_time); err != nil {
			return err.Error()
		}
		householdsCount++
		guestPersonCount += attend_count
		totalGuestMoney += money
	}

	bIsNeedMoneyPack := false
	if householdsCount > MAX_HOUSEHOLD_BOOKED {
		bIsNeedMoneyPack = true
		exceedHousholds := householdsCount - MAX_HOUSEHOLD_BOOKED
		moneyTotalAdded += (exceedHousholds * MONEY_PACK_FOR_BACK_VAL)
		info += fmt.Sprintf("回礼红包个数: %d", exceedHousholds)
	}
	if guestPersonCount > MAX_GUEST_PERSON_BOOKED {
		bIsNeedMoneyPack = true
		exceedGuests := guestPersonCount - MAX_GUEST_PERSON_BOOKED
		moneyTotalAdded += (exceedGuests * MONEY_PACK_FOR_TEA_VAL)
		guestAdded = exceedGuests
		tableAdded = exceedGuests / 10
		if guestAdded%10 >= 5 {
			tableAdded++
		}
		info += fmt.Sprintf("\n茶红包个数(赶席人数): %d", guestAdded)
		info += fmt.Sprintf("\n赶席桌数: %d", tableAdded)
	}
	if bIsNeedMoneyPack {
		info += fmt.Sprintf("\n红包总金额： %d", moneyTotalAdded)
		info += fmt.Sprintf("\n户数：%d  人数：%d  总礼金：%d ￥", householdsCount, guestPersonCount, totalGuestMoney)
		info += fmt.Sprintf("\n打印时间: %s", json_def.TimeStampToString(time.Now().Unix(), 0))
	}
	return
}

func GuestMoneyContribution() *json_def.QueryStatAck {
	var (
		rows                    *sql.Rows
		resp                    error
		householdsCount         int                    //户数
		guestPersonCount        int                    //人数
		totalGuestMoney         int                    //总礼金
		ret                     *json_def.QueryStatAck = &json_def.QueryStatAck{}
		name, entry_time, sql   string                 = "", " ", "select * from t_guest_money"
		id, money, attend_count int

		lastNameMoneyMap  map[string]int = make(map[string]int)
		lastNamePersonMap map[string]int = make(map[string]int)
	)

	if rows, resp = mysql.Cmd(mysql.MYSQL_CMD_QUERY, sql); resp != nil {
		Glog.Error("SubTotalInfo failed: sql=%s err=%s", sql, resp.Error())
		return nil
	}

	for rows.Next() {
		if err := rows.Scan(&id, &name, &money, &attend_count, &entry_time); err != nil {
			return nil
		}

		lastName := string(([]rune(name))[0:1])
		if _, bHas := lastNameMoneyMap[lastName]; !bHas {
			lastNameMoneyMap[lastName], lastNamePersonMap[lastName] = money, 1
		} else {
			lastNameMoneyMap[lastName] += money
			lastNamePersonMap[lastName] += 1
		}
		householdsCount++
		guestPersonCount += attend_count
		totalGuestMoney += money
	}
	if householdsCount != 0 {
		ret.GuestPersonNum = guestPersonCount
		ret.TotalGuestMoney = totalGuestMoney
		ret.HouseHoldNum = householdsCount

		var rankList json_def.RankingList
		rankList.Init(lastNameMoneyMap, lastNamePersonMap)
		sort.Sort(rankList)
		ret.RankList = make([]*json_def.LastName_Money, rankList.Len())
		copy(ret.RankList, rankList)
	}
	return ret
}
