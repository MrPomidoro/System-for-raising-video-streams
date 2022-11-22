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

func (rtsp *rtspRepository) GetRtsp() (map[string]interface{}, error) {
	var item interface{}
	var res map[string]interface{}

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(rtspsimpleserver.URLGetConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		return res, fmt.Errorf("cannot received response from rts-simple-server: %v", err)
	}
	logger.LogDebug(rtsp.log, "Received response from rtsp-simple-server")
	// Отложенное закрытие тела ответа
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	logger.LogDebug(rtsp.log, "Success read body")

	err = json.Unmarshal(body, &item)
	if err != nil {
		return res, fmt.Errorf("cannot unmarshal response: %v", err)
	}
	logger.LogDebug(rtsp.log, "Success unmarshal body")

	res = item.(map[string]interface{})
	return res, nil
}

func (rtsp *rtspRepository) PostAddRTSP(camDB refreshstream.RefreshStream) error {

	// Парсинг поля RunOnReady
	runOnReady := fmt.Sprintf(rtsp.cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)

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
	URLPostAdd := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port, "add", camDB.Stream.String)

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
	URLPostRemove := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port, "remove", camRTSP)

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

	protocol := camDB.Protocol.String

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`{
			"sourceProtocol": "%s",
			"runOnReadRestart": true,
			"runOnReady": "%s"
	}`, protocol, conf.RunOnReady))

	// Парсинг URL
	URLPostEdit := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Server_Host, rtsp.cfg.Server_Port, "edit", camDB.Stream.String)

	// Запрос
	response, err := http.Post(URLPostEdit, "application/json", bytes.NewBuffer(postJson))
	if err != nil {
		return fmt.Errorf("cannot complete post request for edit config: %v", err)
	}
	defer response.Body.Close()

	return nil
}
