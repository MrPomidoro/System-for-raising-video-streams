package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/sirupsen/logrus"
)

type rtspRepository struct {
	cfg *config.Config
	log *logrus.Logger
}

func NewRTSPRepository(cfg *config.Config, log *logrus.Logger) *rtspRepository {
	return &rtspRepository{
		cfg: cfg,
		log: log,
	}
}

func (rtsp *rtspRepository) GetRtsp() map[string]interface{} {
	var item interface{}
	var res map[string]interface{}

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(rtspsimpleserver.URLGetConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		logger.LogError(rtsp.log, fmt.Sprintf("cannot received response from rtspRepository: %v", err))
		return res
	}
	logger.LogDebug(rtsp.log, "Received response from rtsp")
	// Отложенное закрытие тела ответа
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(rtsp.log, err)
		return res
	}
	logger.LogDebug(rtsp.log, "Success read body")

	err = json.Unmarshal(body, &item)
	if err != nil {
		logger.LogError(rtsp.log, fmt.Sprintf("cannot unmarshal response: %v", err))
		return res
	}
	logger.LogDebug(rtsp.log, "Success unmarshal body")

	res = item.(map[string]interface{})
	return res
}

func (rtsp *rtspRepository) PostAddRTSP(camDB refreshstream.RefreshStream) error {

	// Парсинг поля RunOnReady
	runOnReady := fmt.Sprintf(rtspsimpleserver.RunOnReadyConst, rtsp.cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)

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
			"sourceProtocol": "%s",
			"sourceOnDemandStartTimeout": "10s",
			"sourceOnDemandCloseAfter": "10s",
			"readUser": "",
			"readPass": "",
			"runOnDemandStartTimeout": "5s",
			"runOnDemandCloseAfter": "5s",
			"runOnReady": "%s",
			"runOnReadyRestart": true,
			"runOnReadRestart": false
	}`, protocol, runOnReady))

	// Парсинг URL
	URLPostAdd := fmt.Sprintf(rtspsimpleserver.URLPostAddConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port, camDB.Stream.String)

	// Запрос
	response, err := http.Post(URLPostAdd, "application/json; charset=UTF-8", bytes.NewBuffer(postJson))
	if err != nil {
		return fmt.Errorf("cannot complete post request for add config: %v", err)
	}
	defer response.Body.Close()

	return nil
}

func (rtsp *rtspRepository) PostRemoveRTSP(camRTSP string) error {
	// Парсинг URL
	URLPostRemove := fmt.Sprintf(rtspsimpleserver.URLPostRemoveConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port, camRTSP)

	var buf []byte
	// Запрос
	response, err := http.Post(URLPostRemove, "application/json; charset=UTF-8", bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("cannot complete post request for remove config: %v", err)
	}
	defer response.Body.Close()

	return nil
}

func (rtsp *rtspRepository) PostEditRTSP(camDB refreshstream.RefreshStream, conf rtspsimpleserver.Conf) error {

	// Парсинг поля RunOnReady
	// runOnReady := fmt.Sprintf(RunOnReadyConst, cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)

	// Парсинг логина и пароля
	// (не получается занести их в соответствующие поля, как и ip)
	// var login, pass string
	// logPass := strings.Split(camDB.Auth.String, ":")
	// if len(logPass) == 2 {
	// 	login, pass = logPass[0], logPass[1]
	// }

	protocol := camDB.Protocol.String

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`{
			"sourceProtocol": "%s",
			"runOnReadRestart": false,
			"runOnReady": "%s"
	}`, protocol, conf.RunOnReady))

	// Парсинг URL
	URLPostEdit := fmt.Sprintf(rtspsimpleserver.URLPostEditConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port, camDB.Stream.String)

	// Запрос
	response, err := http.Post(URLPostEdit, "application/json", bytes.NewBuffer(postJson))
	if err != nil {
		return fmt.Errorf("cannot complete post request for edit config: %v", err)
	}
	defer response.Body.Close()

	return nil
}

// func (rtsp *rtspRepository) PostSetRTSP(camDB refreshstream.RefreshStream) {
// 	logLevel := rtsp.cfg.LogLevel
// 	logDestinations := "file, os.Stdout"
// 	logFile := rtsp.cfg.LogFile
// 	sourceProtocol := camDB.Protocol.String

// 	fmt.Printf("logDestinations=%s, logLevel=%s, logFile=%s,sourceProtocol=%s\n", logDestinations, logLevel, logFile, sourceProtocol)
// }
