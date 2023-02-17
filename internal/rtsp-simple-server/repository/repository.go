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
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"go.uber.org/zap"
)

type rtspRepository struct {
	cfg *config.Config
	log *zap.Logger
	err *ce.Error
}

func NewRTSPRepository(cfg *config.Config, log *zap.Logger) *rtspRepository {
	return &rtspRepository{
		cfg: cfg,
		log: log,
		err: ce.NewError(ce.ErrorLevel, "50.4.2", "error at rtsp operation level"),
	}
}

// GetRtsp отправляет GET запрос на получение данных
func (rtsp *rtspRepository) GetRtsp() (map[string]interface{}, error) {
	var item interface{}
	var res map[string]interface{}

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(rtspsimpleserver.URLGetConst, rtsp.cfg.Url)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	// Отложенное закрытие тела ответа
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	rtsp.log.Debug("Received response from rtsp-simple-server")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	rtsp.log.Debug("Success read body")

	err = json.Unmarshal(body, &item)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	rtsp.log.Debug("Success unmarshal body")

	res = item.(map[string]interface{})
	return res, nil
}

// PostAddRTSP отправляет POST запрос на добавление потока
func (rtsp *rtspRepository) PostAddRTSP(camDB refreshstream.RefreshStream) error {

	// Парсинг поля RunOnReady
	var runOnReady string
	if rtsp.cfg.Run != "" {
		runOnReady = fmt.Sprintf(rtsp.cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)
	} else {
		runOnReady = ""
	}

	// Поле протокола не должно быть пустым
	// по умолчанию - tcp
	var protocol = camDB.Protocol.String
	if protocol == "" {
		protocol = "tcp"
	}

	source := fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream.String)

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`{
			"source": "%s",
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
	}`, source, protocol, runOnReady))

	// Парсинг URL
	URLPostAdd := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Url, "add", camDB.Stream.String)

	// Запрос
	resp, err := http.Post(URLPostAdd, "application/json; charset=UTF-8", bytes.NewBuffer(postJson))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return rtsp.err.SetError(err)
	}

	return nil
}

// PostRemoveRTSP отправляет POST запрос на удаление потока
func (rtsp *rtspRepository) PostRemoveRTSP(camRTSP string) error {
	// Парсинг URL
	URLPostRemove := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Url, "remove", camRTSP)

	var buf []byte
	// Запрос
	resp, err := http.Post(URLPostRemove, "application/json; charset=UTF-8", bytes.NewBuffer(buf))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return rtsp.err.SetError(err)
	}

	return nil
}

// PostEditRTSP отправляет POST запрос на изменение потока
func (rtsp *rtspRepository) PostEditRTSP(camDB refreshstream.RefreshStream, conf rtspsimpleserver.Conf) error {

	var protocol = camDB.Protocol.String
	if protocol == "" && conf.SourceProtocol == "" {
		protocol = "tcp"
	}

	var runOnReady = conf.RunOnReady
	if rtsp.cfg.Run != "" {
		runOnReady = fmt.Sprintf(rtsp.cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)
	} else {
		runOnReady = ""
	}

	var source = conf.Source
	if source == "" {
		source = fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream.String)
	}

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`{
			"source": "%s",
			"sourceProtocol": "%s",
			"runOnReadRestart": false,
			"runOnReady": "%s"
	}`, source, protocol, runOnReady))

	// Парсинг URL
	URLPostEdit := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Url, "edit", camDB.Stream.String)

	// Запрос
	resp, err := http.Post(URLPostEdit, "application/json", bytes.NewBuffer(postJson))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return rtsp.err.SetError(err)
	}

	return nil
}
