package main

import (
	// "bytes"
	// "encoding/json"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	// "time"

	// "github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	// "net/http"
	// "time"
)

type Notification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:receiver`
	GroupLabels       map[string]string `json:groupLabels`
	CommonLabels      map[string]string `json:commonLabels`
	CommonAnnotations map[string]string `json:commonAnnotations`
	ExternalURL       string            `json:externalURL`
}
type res struct {
	Result string `json:"result"`
}
type AlertInfo struct {
	Result []Alert `json:"result"`
}
type Alert struct {
	Description string   `json:"description"`
	Hosts       []*Hosts `json:"hosts"`
}
type Hosts struct {
	Host string `json:"host"`
}
type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingTalkMarkdown struct {
	MsgType  string    `json:"msgtype"`
	At       *At       `json:at`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func main() {
	s := res{}
	gout.POST("http://ali-zabbix.ops.yzs.io//zabbix/api_jsonrpc.php").
		//打开debug模式
		Debug(false).
		//设置json到请求body
		SetJSON(
			gout.H{
				"jsonrpc": "2.0",
				"method":  "user.login",
				"params": gout.H{
					"user":     "Admin",
					"password": "xxx",
				},
				"id":   1,
				"auth": nil,
			},
		).
		BindJSON(&s).
		Do()

	getalet(s.Result)

}
func getalet(auth string) {
	alertinfo := AlertInfo{}
	gout.POST("http://ali-zabbix.ops.yzs.io//zabbix/api_jsonrpc.php").
		//打开debug模式
		Debug(false).
		//设置json到请求body
		SetJSON(
			gout.H{
				"jsonrpc": "2.0",
				"method":  "trigger.get",
				"params": gout.H{
					"output":            gout.A{"triggerid", "description", "priority"},
					"filter":            gout.H{"value": 1},
					"sortfield":         "priority",
					"sortorder":         "DESC",
					"min_severity":      1,
					"skipDependent":     1,
					"monitored":         1,
					"active":            1,
					"expandDescription": 1,
					"selectHosts":       gout.A{"host"},
					"selectGroups":      gout.A{"name"},

					"only_true": 1,
				},
				"id":   1,
				"auth": auth,
			},
		).
		BindBody(&alertinfo).
		Do()
	Send(alertinfo, "https://oapi.dingtalk.com/robot/send?access_token=dxxxxx")
}
func TransformToMarkdown(notification AlertInfo) (markdown *DingTalkMarkdown, err error) {

	Alertinfo := notification.Result

	var buffer bytes.Buffer

	for _, alert := range Alertinfo {
		annotations := alert.Description
		for _, host := range alert.Hosts {
			res := strings.Contains(host.Host, "asr")
			if res {
				str, _ := json.Marshal(alert)
				fmt.Println(str)
				buffer.WriteString(fmt.Sprintf("##### %s引擎异常\n > %s\n", host.Host, annotations))
				buffer.WriteString(fmt.Sprintf("\n> 主机名：%s\n", host.Host))
				buffer.WriteString(fmt.Sprintf("\n> 开始时间：%s\n", time.Now().Format("2006-01-02 15:04:05")))
			}
		}

	}

	markdown = &DingTalkMarkdown{
		MsgType: "markdown",
		Markdown: &Markdown{
			Title: fmt.Sprintf("报警提醒"),
			Text:  buffer.String(),
		},
		At: &At{
			IsAtAll: false,
		},
	}
	return
}

func Send(alertinfo AlertInfo, dingtalkRobot string) (err error) {

	markdown, err := TransformToMarkdown(alertinfo)

	if err != nil {
		return
	}

	data, err := json.Marshal(markdown)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		"POST",
		dingtalkRobot,
		bytes.NewBuffer(data))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return
}
