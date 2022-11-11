package rtsp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

func GetRtsp(cfg *config.Config) map[string]interface{} {
	log := logger.NewLog(cfg.LogLevel)
	var item interface{}
	var res map[string]interface{}

	resp, err := http.Get("http://10.100.100.30:9997/v1/paths/list")
	if err != nil {
		logger.LogError(log, fmt.Sprintf("cannot отправить сраный запрос на rtsp: %v", err))
		return res
	}
	// logger.LogDebug(log, fmt.Sprintf("response:\n%+v", resp.Body))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(log, err)
		return res
	}
	logger.LogDebug(log, "Success read body")

	// fmt.Println(string(body))

	err = json.Unmarshal(body, &item)
	if err != nil {
		logger.LogError(log, fmt.Sprintf("cannot unmarshal response: %v", err))
		return res
	} else {
		logger.LogDebug(log, "Success unmarshal body")
	}

	res = item.(map[string]interface{})
	return res
}
