package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/transcode"
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
func (rtsp *rtspRepository) GetRtsp(ctx context.Context) (map[string]rtspsimpleserver.SConf, ce.IError) {

	// defer close(dataRTSPchan)
	res := make(map[string]rtspsimpleserver.SConf)

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(rtsp.cfg.Url + rtsp.cfg.Rtsp.Api.UrlGet)
	// URLGet := fmt.Sprintf(rtspsimpleserver.URLGetConst, rtsp.cfg.Url)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLGet)
	// Get запрос и обработка ошибки
	resp, err := http.Get(URLGet)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	// Закрытие тела ответа
	defer resp.Body.Close()
	rtsp.log.Debug("Received response from rtsp-simple-server")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	rtsp.log.Debug("Success read body")

	var item map[string]interface{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	if len(item) == 0 {
		return res, rtsp.err.SetError(fmt.Errorf("response from rtsp not received"))
	}

	rtsp.log.Debug("Success unmarshal body")

	for _, ress := range item {
		item1 := ress.(map[string]interface{})

		for stream, i := range item1 {
			cam := rtspsimpleserver.SConf{}
			cam.Stream = stream

			fileds := i.(map[string]interface{})
			transcode.Transcode(fileds["conf"], &cam.Conf)
			res[stream] = cam
		}
	}
	return res, nil
}

// PostAddRTSP отправляет POST запрос на добавление потока
func (rtsp *rtspRepository) PostAddRTSP(ctx context.Context, cam rtspsimpleserver.SConf) ce.IError {

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`
	{
		"source": "%s",
		"sourceProtocol": "%s",
		"sourceOnDemandCloseAfter": "10s",
		"sourceOnDemandStartTimeout": "10s",
		"readPass": "",
		"readUser": "",
		"runOnDemandCloseAfter": "10s",
		"runOnDemandStartTimeout": "10s",
		"runOnReadyRestart": true,
		"runOnReady": "%s",
		"runOnReadRestart": false
	}`, cam.Conf.Source, cam.Conf.SourceProtocol, cam.Conf.RunOnReady))

	// Парсинг URL
	URLPostAdd := fmt.Sprintf(rtsp.cfg.Url + rtsp.cfg.Rtsp.Api.UrlAdd + cam.Stream)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostAdd)

	if ctx.Err() != nil {
		return rtsp.err.SetError(ctx.Err())
	}
	// Запрос
	resp, err := http.Post(URLPostAdd, "application/json; charset=UTF-8", bytes.NewBuffer(postJson))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	defer resp.Body.Close()

	return nil
}

// PostRemoveRTSP отправляет POST запрос на удаление потока
func (rtsp *rtspRepository) PostRemoveRTSP(ctx context.Context, camRTSP rtspsimpleserver.SConf) ce.IError {
	// Парсинг URL
	URLPostRemove := fmt.Sprintf(rtsp.cfg.Url + rtsp.cfg.Rtsp.Api.UrlRemove + camRTSP.Stream)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostRemove)

	var buf []byte

	if ctx.Err() != nil {
		return rtsp.err.SetError(ctx.Err())
	}
	// Запрос
	resp, err := http.Post(URLPostRemove, "application/json; charset=UTF-8", bytes.NewBuffer(buf))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	defer resp.Body.Close()
	return nil
}

// PostEditRTSP отправляет POST запрос на изменение потока
func (rtsp *rtspRepository) PostEditRTSP(ctx context.Context, cam rtspsimpleserver.SConf) ce.IError {

	// Формирование джейсона для отправки
	postJson := []byte(fmt.Sprintf(`
	{
		"source": "%s",
		"sourceProtocol": "%s",
		"runOnReady": "%s",
		"runOnReadRestart": false
	}`, cam.Conf.Source, cam.Conf.SourceProtocol, cam.Conf.RunOnReady))

	// Парсинг URL
	URLPostEdit := fmt.Sprintf(rtsp.cfg.Url + rtsp.cfg.Rtsp.Api.UrlEdit + cam.Stream)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostEdit)

	if ctx.Err() != nil {
		return rtsp.err.SetError(ctx.Err())
	}
	// Запрос
	resp, err := http.Post(URLPostEdit, "application/json", bytes.NewBuffer(postJson))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	defer resp.Body.Close()

	return nil
}
