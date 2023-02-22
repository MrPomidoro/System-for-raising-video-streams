package repository

import (
	"bytes"
	"context"
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
	err ce.IError
}

func NewRTSPRepository(cfg *config.Config, log *zap.Logger) *rtspRepository {
	return &rtspRepository{
		cfg: cfg,
		log: log,
		err: ce.ErrorRTSP,
	}
}

// GetRtsp отправляет GET запрос на получение данных
func (rtsp *rtspRepository) GetRtsp(ctx context.Context,
	dataRTSPchan chan rtspsimpleserver.SConf) ([]rtspsimpleserver.SConf, ce.IError) {

	defer close(dataRTSPchan)

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(rtspsimpleserver.URLGetConst, rtsp.cfg.Url)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLGet)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		return nil, rtsp.err.SetError(err)
	}
	// Закрытие тела ответа
	defer resp.Body.Close()
	rtsp.log.Debug("Received response from rtsp-simple-server")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rtsp.err.SetError(err)
	}
	rtsp.log.Debug("Success read body")

	var item map[string]interface{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return nil, rtsp.err.SetError(err)
	}
	if len(item) == 0 {
		return nil, rtsp.err.SetError(fmt.Errorf("response from rtsp not received"))
	}

	rtsp.log.Debug("Success unmarshal body")

	var res []rtspsimpleserver.SConf
	for _, ress := range item {
		ress1 := ress.(map[string]interface{})
		for stream, i := range ress1 {
			sconf := rtspsimpleserver.SConf{}
			sconf.Stream = stream
			im := i.(map[string]interface{})
			for field, j := range im {
				if field == "conf" {
					transcode(j, &sconf.Conf)
					break
				}
			}

			res = append(res, sconf)

			select {
			case dataRTSPchan <- sconf:
			case <-ctx.Done():
				return res, nil
			}
		}
	}
	// fmt.Printf("%+v\n\n", sConfs)
	return res, nil
}

func transcode(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

// PostAddRTSP отправляет POST запрос на добавление потока
func (rtsp *rtspRepository) PostAddRTSP(camDB refreshstream.RefreshStream) ce.IError {

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
	postJson := []byte(fmt.Sprintf(`
	{
		"sourceProtocol": "%s",
		"source": "%s",
		"sourceOnDemandCloseAfter": "10s",
		"sourceOnDemandStartTimeout": "10s",
		"readPass": "",
		"readUser": "",
		"runOnDemandCloseAfter": "5s",
		"runOnDemandStartTimeout": "5s",
		"runOnReadyRestart": true,
		"runOnReady": "%s",
		"runOnReadRestart": false
	}`, source, protocol, runOnReady))

	// Парсинг URL
	URLPostAdd := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Url, "add", camDB.Stream.String)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostAdd)

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
func (rtsp *rtspRepository) PostRemoveRTSP(camRTSP string) ce.IError {
	// Парсинг URL
	URLPostRemove := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Url, "remove", camRTSP)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostRemove)

	var buf []byte

	// Запрос
	resp, err := http.Post(URLPostRemove, "application/json; charset=UTF-8", bytes.NewBuffer(buf))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	resp.Body.Close()
	return nil
}

// PostEditRTSP отправляет POST запрос на изменение потока
func (rtsp *rtspRepository) PostEditRTSP(camDB refreshstream.RefreshStream, sconf rtspsimpleserver.SConf) ce.IError {

	var protocol = camDB.Protocol.String
	if protocol == "" && sconf.Conf.SourceProtocol == "" {
		protocol = "tcp"
	}

	var runOnReady = sconf.Conf.RunOnReady
	if rtsp.cfg.Run != "" {
		runOnReady = fmt.Sprintf(rtsp.cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)
	} else {
		runOnReady = ""
	}

	var source = sconf.Conf.Source
	if source == "" {
		source = fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream.String)
	}

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`
	{
		"sourceProtocol": "%s",
		"source": "%s",
		"runOnReady": "%s"
		"runOnReadRestart": false,
	}`, source, protocol, runOnReady))

	// Парсинг URL
	URLPostEdit := fmt.Sprintf(rtspsimpleserver.URLPostConst, rtsp.cfg.Url, "edit", camDB.Stream.String)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostEdit)

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
