package json_def

import (
	"strings"
	"time"
)

type QueryGuestAck struct {
	Id          int
	Name        string
	Money       int
	AttendCount int
	EntryTime   string
}

type GuestRankHelper []*QueryGuestAck

func (g GuestRankHelper) Len() int {
	return len(g)
}

func (g GuestRankHelper) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GuestRankHelper) Less(i, j int) bool {
	return g[i].Id > g[j].Id
	iTimeStamp := DateTimeStringToStamp(g[i].EntryTime)
	jTimeStamp := DateTimeStringToStamp(g[j].EntryTime)
	if iTimeStamp != jTimeStamp {
		return iTimeStamp > jTimeStamp
	} else if g[i].Money != g[j].Money {
		return g[i].Money > g[j].Money
	}
	return false
}

type QueryStatAck struct {
	HouseHoldNum    int
	GuestPersonNum  int
	TotalGuestMoney int
	RankList        []*LastName_Money
}

type LastName_Money struct {
	LastName    string
	TotalMoney  int
	PersonCount int
}

type RankingList []*LastName_Money

func (rl *RankingList) Init(moneyTable, personTable map[string]int) {
	for lastName, totalMoney := range moneyTable {
		node := &LastName_Money{lastName, totalMoney, 0}
		if val, has := personTable[lastName]; has {
			node.PersonCount = val
		}
		*rl = append(*rl, node)
	}
}

func (rl RankingList) Len() int {
	return len(rl)
}

func (rl RankingList) Less(i, j int) bool {
	if rl[i].TotalMoney != rl[j].TotalMoney {
		return rl[i].TotalMoney > rl[j].TotalMoney
	} else if rl[i].PersonCount != rl[j].PersonCount {
		return rl[i].PersonCount < rl[j].PersonCount
	}
	return strings.Compare(rl[i].LastName, rl[j].LastName) < 0
}

func (rl RankingList) Swap(i, j int) {
	rl[i], rl[j] = rl[j], rl[i]
}

func TimeStampToString(secs, nanoSecs int64) string {
	layout := "2006-01-02 15:04:06"
	return time.Unix(secs, nanoSecs).Format(layout)
}

func DateTimeStringToStamp(dateTimeStr string) int64 {
	locale, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	layout := "2006-01-02 15:04:06"
	theTime, err := time.ParseInLocation(layout, dateTimeStr, locale)
	if err != nil {
		panic(err)
	}
	return theTime.Unix()
}

type CSCommonAck struct {
}

func PackJsonAck(code int, msg string, dataAck interface{}) interface{} {
	var commAck struct {
		ErrCode int
		ErrMsg  string
		Data    interface{}
	}
	commAck.ErrCode, commAck.ErrMsg, commAck.Data = code, msg, dataAck
	return commAck
}
