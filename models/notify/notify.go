package notify

import (
	. "birthday_server/models/log"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	T_NOTIFY_MSG = iota
	T_NOTIFY_STAT
)
const (
	get_token_http_format_url   = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	send_notify_http_format_url = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	global_corpid               = "wwfce4639dddb21018"
	msg_notify_app_agentid      = 1000003
	msg_notify_app_scretid      = "KRjAUS2coRBR1jiFb6ta_FaR5yFuhrpilaGdmy18MDU"
	stat_notify_app_agentid     = 1000004
	stat_notify_app_scretid     = "eIw1Uz--jPl50jox2NKSqu_KevvYORxDWQn8z2LzPGI"
)

type NotifyApp struct {
	appid      int    //is same with agent_id
	Token      string //token for this app
	flush_time int64  //记录最近一次获取token成功的时间戳
}

var (
	notifyScretMap map[int]string = map[int]string{
		msg_notify_app_agentid:  msg_notify_app_scretid,
		stat_notify_app_agentid: stat_notify_app_scretid,
	}
	GMsgNotify, GStatNotify *NotifyApp
)

func init() {
	GMsgNotify = &NotifyApp{msg_notify_app_agentid, "", 0}
	GStatNotify = &NotifyApp{stat_notify_app_agentid, "", 0}
	fmt.Printf("notify Component started...")
}

func (a *NotifyApp) GetNotifyAccessToken() (token string, err error) {
	secretid := ""
	if val, has := notifyScretMap[a.appid]; !has {
		Glog.Error("the type of msg can't be surported now: %d", a.appid)
		return "", fmt.Errorf("the type of msg can't be surported now: %d", a.appid)
	} else {
		secretid = val
	}
	now := time.Now().Unix()
	if a.Token == "" || now-a.flush_time > 7200 {
		url := fmt.Sprintf(get_token_http_format_url, global_corpid, secretid)
		var resp *http.Response
		Glog.Debug("get token by: %s", url)
		resp, err = http.Get(url)
		if err != nil {
			Glog.Error("get token failed: err=%s", err.Error())
			return "", err
		}
		if resp.StatusCode != 200 {
			Glog.Error("some error occured: errcode=%d errmsg=%s", resp.StatusCode, resp.Status)
			return "", fmt.Errorf(resp.Status)
		}

		var server_node_resp struct {
			Code   int    `json:errcode`
			Errmsg string `json:"errmsg"`
			Token  string `json:"access_token"`
			Expire int64  `json:"expires_in"`
		}

		buf, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return "", e
		}

		if err = json.Unmarshal(buf, &server_node_resp); err != nil {
			Glog.Error("make sure same data formate between caller and server")
			return "", err
		}

		a.Token = server_node_resp.Token
		a.flush_time = now
	}
	return a.Token, nil
}

type NotifyRequest struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

func (n *NotifyApp) NotifyEveryOne(content string) error {
	token, err := n.GetNotifyAccessToken()
	if err != nil {
		Glog.Error("get notify token failed: err=%s", err.Error())
		return err
	}
	send_url := fmt.Sprintf(send_notify_http_format_url, token)
	var (
		resp *http.Response
		body *bytes.Buffer
	)
	body = bytes.NewBuffer([]byte(n.make_send_text(content)))
	resp, err = http.Post(send_url, "application/json", body)
	if err != nil {
		Glog.Error("send notify failed: err=%s", err.Error())
		return err
	}
	var resp_server_node struct {
		Errcode int    `json:"errcode`
		Errmsg  string `json:"errmsg"`
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(buf, &resp_server_node); err != nil {
		Glog.Error("json uncode body failed: %s", err.Error())
		return err
	}
	if resp_server_node.Errcode != 0 {
		return fmt.Errorf(string(buf))
	}
	Glog.Debug("notify ok: msg=%s", content)
	return nil
}

func (n *NotifyApp) make_send_text(content string) (msg []byte) {
	req := &NotifyRequest{
		Touser:  "@all",
		Toparty: "@all",
		Totag:   "@all",
		Msgtype: "text",
	}
	textMap := make(map[string]string)
	textMap["content"] = content
	req.Text = textMap
	req.Agentid = n.appid
	req.Safe = 0
	if buff, err := json.Marshal(req); err != nil {
		Glog.Error("pack req data failed: %s", err.Error())
		return
	} else {
		msg = buff
	}
	return
}

func TestGetToken() {
	token, err := GMsgNotify.GetNotifyAccessToken()
	Glog.Debug("token: %s err: %+v", token, err)
}

func TestSendNotify() {
	if err := GMsgNotify.NotifyEveryOne("hello,everyone!"); err != nil {
		Glog.Debug("send notify to everyone failed...")
	}
}
