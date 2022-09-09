package monitor

import (
	"bytes"
	"context"
	"encoding/json"

	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

type Wsres struct {
	Result  Result `json:"result"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type Result struct {
	Score float64 `json: "score"`
}

func CnWsMonitor(area string, url string) (code int, res string) {
	ss := `{"appkey":"nulk63ngwrfc4xrfp7elud6aeout2ueubeunomis","audioFormat":"mp3","displayText":"你好","eof":"guo-test-end","EvalType":"paragraph","scoreCoefficient":"1","Language":"cn"}`

	audio, err := ioutil.ReadFile("./scripts/nihao.mp3")
	if err != nil {
		log.Printf("%s read audio file err %v\n", area, err)
		return -1, "read audio file err"
	}

	result, err := GetEvalWebsocketResult(url, ss, audio)
	if err != nil {
		log.Printf("%s eval websocket err %v\n", area, err)
		return -1, err.Error()
	} else {
		log.Println(result)
		wsres := Wsres{}
		err := json.Unmarshal([]byte(result), &wsres)
		if err != nil {
			log.Printf("%s %v\n", area, err)
			return -1, err.Error()
		} else {
			if wsres.Errcode == 0 {
				return 200, "success"
			} else {
				return wsres.Errcode, wsres.Errmsg
			}
		}
	}
}
func EnWsMonitor(area string, url string) (code int, res string) {
	ss := `{"appkey":"nulk63ngwrfc4xrfp7elud6aeout2ueubeunomis","audioFormat":"pcm","displayText":"hello world","eof":"guo-test-end","EvalType":"paragraph","scoreCoefficient":"1","Language":"en"}`

	audio, err := ioutil.ReadFile("./scripts/helloworld.wav")
	if err != nil {
		log.Printf("%s read audio file err %v\n", area, err)
		return -1, "read audio file err"
	}

	result, err := GetEvalWebsocketResult(url, ss, audio)
	if err != nil {
		log.Printf("%s eval websocket err %v\n", area, err)
		return -1, err.Error()
	} else {
		log.Println(result)
		wsres := Wsres{}
		err := json.Unmarshal([]byte(result), &wsres)
		if err != nil {
			log.Printf("%s %v\n", area, err)
			return -1, err.Error()
		} else {
			if wsres.Errcode == 0 {
				return 200, "success"
			} else {
				return wsres.Errcode, wsres.Errmsg
			}
		}
	}
}
func GetEvalWebsocketResult(ip string, text string, audio []byte) (result string, err error) {
	url := fmt.Sprintf("ws://%s/ws/eval/", ip)
	conn, _, err := websocket.DefaultDialer.DialContext(context.Background(), url, nil)
	if err != nil {
		return result, fmt.Errorf("websocket dialContext url %s err %v", url, err)
	}
	_ = conn.WriteMessage(websocket.TextMessage, []byte(text))
	tmpData := bytes.NewReader(audio)
	for {
		buf := make([]byte, 8000)
		n, err := io.ReadFull(tmpData, buf)
		if err == io.EOF {
			break
		}
		err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
		if err != nil {
			return result, fmt.Errorf("send websocket audio err %v", err)
		}
	}
	eof := gjson.Get(text, "eof").String()
	_ = conn.WriteMessage(websocket.TextMessage, []byte(eof))
	_, out, err := conn.ReadMessage()
	if err != nil {
		return result, fmt.Errorf("websocket read message err %v", err)
	}
	result = string(out)
	return result, nil
}
