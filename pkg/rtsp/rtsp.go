package rtsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

func GetRtsp(cfg *config.Config) map[string]interface{} {
	log := logger.NewLog(cfg.LogLevel)
	var item interface{}
	var res map[string]interface{}

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(URLGetConst, cfg.Server_Host, cfg.Server_Port)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		logger.LogError(log, fmt.Sprintf("cannot received response from rtsp: %v", err))
		return res
	}
	logger.LogInfo(log, "Received response from rtsp")
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
	}
	logger.LogDebug(log, "Success unmarshal body")

	res = item.(map[string]interface{})
	return res
}

func PostAddRTSP(camDB refreshstream.RefreshStream, cfg *config.Config) error {

	// Парсинг поля RunOnReady
	runOnReady := fmt.Sprintf(RunOnReadyConst, cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)

	// Парсинг логина и пароля
	// (не получается занести их в соответствующие поля, как и ip)
	// var login, pass string
	// logPass := strings.Split(camDB.Auth.String, ":")
	// if len(logPass) == 2 {
	// 	login, pass = logPass[0], logPass[1]
	// }

	// Поле протокола не должно быть пустым
	// по умолчанию - tcp
	var protocol string = camDB.Protocol.String
	if protocol == "" {
		protocol = "tcp"
	}

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`{
			"source": "",
			"sourceProtocol": "%s",
			"sourceAnyPortEnable": false,
			"sourceFingerprint": "",
			"sourceOnDemand": false,
			"sourceOnDemandStartTimeout": "10s",
			"sourceOnDemandCloseAfter": "10s",
			"sourceRedirect": "",
			"disablePublisherOverride": false,
			"fallback": "",
			"publishUser": "",
			"publishPass": "",
			"readUser": "",
			"readPass": "",
			"runOnInit": "",
			"runOnInitRestart": false,
			"runOnDemand": "",
			"runOnDemandRestart": true,
			"runOnDemandStartTimeout": "5s",
			"runOnDemandCloseAfter": "5s",
			"runOnReady": "%s",
			"runOnReadyRestart": false,
			"runOnRead": "",
			"runOnReadRestart": false
	}`, protocol, runOnReady))

	// Парсинг URL
	URLPostAdd := fmt.Sprintf(URLPostAddConst, cfg.Server_Host, cfg.Server_Port, camDB.Stream.String)

	// Запрос
	response, err := http.Post(URLPostAdd, "application/json; charset=UTF-8", bytes.NewBuffer(postJson))
	if err != nil {
		return fmt.Errorf("cannot complete post request for add config: %v", err)
	}
	defer response.Body.Close()

	return nil
}

func PostRemoveRTSP(camRTSP string, cfg *config.Config) error {
	// Парсинг URL
	URLPostRemove := fmt.Sprintf(URLPostRemoveConst, cfg.Server_Host, cfg.Server_Port, camRTSP)

	var buf []byte
	// Запрос
	response, err := http.Post(URLPostRemove, "application/json; charset=UTF-8", bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("cannot complete post request for remove config: %v", err)
	}
	defer response.Body.Close()

	return nil
}
