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
	logStC := logger.NewLogStatCode(cfg.LogLevel)
	log := logger.NewLog(cfg.LogLevel)
	var item interface{}
	var res map[string]interface{}

	resp, err := http.Get("http://10.100.100.228:9997/v1/paths/list")
	if err != nil {
		logger.LogErrorStatusCode(logStC, fmt.Sprintf("cannot to send request to rtsp: %v", err), "Get", "500")
		return res
	}
	logger.LogInfoStatusCode(logStC, "Success send Get request to rtsp", "Get", "200")
	// logger.LogDebug(log, fmt.Sprintf("response:\n%+v", resp.Body))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(log, err)
		return res
	}
	logger.LogDebug(log, "Success read body")

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
