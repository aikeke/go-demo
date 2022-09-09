package monitor

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/guonaihong/gout"
)

func EnHttpMonitor(area, url string) (code int, res string) {
	s := ""
	code = 0
	t := testForm{
		Mode:  "B",
		Text:  "hello world",
		Voice: "./scripts/helloworld.wav",
	}

	err := gout.POST(url).SetHeader(gout.H{"session-id": "uuidgen", "appkey": "xuqo7pqagqx5gvdbqyfybrusfosbbkjjtfvsr5qx"}).SetForm(&t).BindBody(&s).Code(&code).Do()

	if err != nil {
		log.Println(area + "英文http评测,网络异常\n" + err.Error())
		return -1, "网络异常"
	}
	var tempMap map[string]interface{}
	err = json.Unmarshal([]byte(s), &tempMap)

	if err != nil {
		return -1, s
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
