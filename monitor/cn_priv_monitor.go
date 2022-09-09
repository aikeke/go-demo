package monitor

import (
	"encoding/json"

	"log"
	"os/exec"
)

func CnPrivMonitor(area string, url string) (code int, result string) {
	cmd := `./scripts/priv_tst.sh cn ` + url
	res, _ := exec.Command("sh", "-c", cmd).Output()
	privres := PrivRes{}
	err := json.Unmarshal(res, &privres)
	if err != nil {
		log.Println("json解析失败", err)
		return -1, string(res)
	}
	if privres.Message == "" {
		return 200, "success"
	} else {
		log.Println(string(res))
		return -2, privres.Message
	}

}
