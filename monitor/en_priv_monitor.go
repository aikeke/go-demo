package monitor

import (
	"context"
	"encoding/json"
	"time"

	"log"
	"os/exec"
)

type PrivRes struct {
	Score   float64 `json:"score"`
	Message string  `json:"message"`
}

func EnPrivMonitor(area string, url string) (code int, result string) {
	ctx, cannel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cannel()
	cmd := `./scripts/priv_tst.sh en ` + url
	res, _ := exec.CommandContext(ctx, "sh", "-c", cmd).Output()

	privres := PrivRes{}
	err := json.Unmarshal(res, &privres)
	if err != nil {
		log.Println("json解析失败", err)
		return -2, string(res)
	}
	if privres.Message == "" {
		return 200, "success"
	} else {
		log.Println(string(res))
		return -1, privres.Message
	}

}
