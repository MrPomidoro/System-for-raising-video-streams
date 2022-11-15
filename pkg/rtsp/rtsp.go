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

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf("http://%s%s/v1/paths/list", cfg.Server_Host, cfg.Server_Port)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		logger.LogErrorStatusCode(logStC, fmt.Sprintf("cannot to send request to rtsp: %v", err), "Get", "500")
		return res
	}
	logger.LogInfoStatusCode(logStC, "Received response from the rtsp", "Get", "200")
	// Отложенное закрытие тела ответа
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
