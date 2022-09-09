package monitor

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/guonaihong/gout"
)

type testForm struct {
	Mode  string `form:"mode"`
	Text  string `form:"text"`
	Voice string `form:"voice" form-file:"true"` //从文件中读取
}

func CnHttpMonitor(area, url string) (code int, res string) {
	s := ""
	code = 0
	t := testForm{
		Mode:  "B",
		Text:  "你好",
		Voice: "./scripts/nihao.mp3",
	}

	err := gout.POST(url).SetHeader(gout.H{"session-id": "uuidgen", "appkey": "xuqo7pqagqx5gvdbqyfybrusfosbbkjjtfvsr5qx"}).SetForm(&t).BindBody(&s).Code(&code).Do()

	if err != nil {
		log.Println(area + "中文http评测,网络异常\n" + err.Error())
		return -1, "网络异常"
	}
	var tempMap map[string]interface{}
	err = json.Unmarshal([]byte(s), &tempMap)

	if err != nil {
		return -1, err.Error()

	}
	scorelist := tempMap["score"]
	switch v := scorelist.(type) {

	case float64:
		if v < 80.0 {
			score := strconv.Itoa(int(v))
			return -2, "评测分数异常，当前分数" + score
		} else {
			return code, "success"
		}
	default:
		return code, s
	}

}
