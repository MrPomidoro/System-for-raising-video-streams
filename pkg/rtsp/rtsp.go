package rtsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

func GetRtsp(cfg *config.Config) map[string]interface{} {
	logStC := logger.NewLogStatCode(cfg.LogLevel)
	log := logger.NewLog(cfg.LogLevel)
	var item interface{}
	var res map[string]interface{}

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(URLGetConst, cfg.Server_Host, cfg.Server_Port)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		logger.LogErrorStatusCode(logStC, fmt.Sprintf("cannot received response from rtsp: %v", err), "Get", "500")
		return res
	}
	logger.LogInfoStatusCode(logStC, "Received response from rtsp", "Get", "200")
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

func PostRTSP(camDB refreshstream.RefreshStream, cfg *config.Config) error {

	// Парсинг поля RunOnReady
	runOnReady := fmt.Sprintf(RunOnReadyConst, cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)

	// Парсинг логина и пароля
	var login, pass string
	logPass := strings.Split(camDB.Auth.String, ":")
	if len(logPass) == 2 {
		login, pass = logPass[0], logPass[1]
	}

	// Заполнение структуры Conf для отправления Post запроса
	postStruct := Conf{RunOnReadRestart: true, ReadIPs: []string{camDB.Ip.String},
		RunOnReady: runOnReady, ReadUser: login, ReadPass: pass, SourceProtocol: camDB.Protocol.String}

	// fmt.Printf("%#v\n\n", postStruct)

	// Формирование мапы
	fieldsMap := make(map[string]interface{})
	fieldsMap["conf"] = postStruct
	fieldsMap["confName"] = camDB.Stream.String
	fieldsMap["source"] = ""
	fieldsMap["sourceReady"] = false
	fieldsMap["tracks"] = []string{}
	fieldsMap["readers"] = []string{}

	postMap := make(map[string]map[string]interface{})
	postMap[camDB.Stream.String] = fieldsMap

	// fmt.Printf("%#v\n\n", postMap)

	// Маршалинг в json
	postJson, err := json.Marshal(postMap)
	if err != nil {
		return fmt.Errorf("cannot marshal structure to json, %v", err)
	}
	fmt.Println(string(postJson))

	// Парсинг URL
	URLPost := fmt.Sprintf(URLPostConst, cfg.Server_Host, cfg.Server_Port, camDB.Stream.String)
	// Запрос
	_, err = http.NewRequest("POST", URLPost, bytes.NewBuffer(postJson))
	if err != nil {
		return fmt.Errorf("cannot complete post request, %v", err)
	}

	return nil
}
