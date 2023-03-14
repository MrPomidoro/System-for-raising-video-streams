package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"go.uber.org/zap"
)

type repository struct {
	cfg    *config.Config
	log    *zap.Logger
	err    ce.IError
	client *http.Client
}

func NewRepository(cfg *config.Config, log *zap.Logger) *repository {
	return &repository{
		cfg:    cfg,
		log:    log,
		err:    ce.ErrorRTSP,
		client: &http.Client{},
	}
}

// GetRtsp отправляет GET запрос на получение данных
func (rtsp *repository) GetRtsp(ctx context.Context) (map[string]rtspsimpleserver.SConf, ce.IError) {

	res := make(map[string]rtspsimpleserver.SConf)

	// Формирование URL для get запроса
	URLGet := fmt.Sprintf(rtsp.cfg.Url + rtsp.cfg.Rtsp.Api.UrlGet)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLGet)

	payload := strings.NewReader(``)

	req, err := http.NewRequest(http.MethodGet, URLGet, payload)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}

	if ctx.Err() != nil {
		return res, rtsp.err.SetError(ctx.Err())
	}

	// Get запрос и обработка ошибки
	resp, err := rtsp.client.Do(req)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	// Закрытие тела ответа
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return res, rtsp.err.SetError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	rtsp.log.Debug("Received response from rtsp-simple-server")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	rtsp.log.Debug("Success read body")

	var item map[string]map[string]rtspsimpleserver.SConf
	err = json.Unmarshal(body, &item)
	if err != nil {
		return res, rtsp.err.SetError(err)
	}
	if len(item) == 0 {
		return res, rtsp.err.SetError(fmt.Errorf("response from rtsp not received"))
	}

	rtsp.log.Debug("Success unmarshal body")

	for _, ress := range item {
		for stream, i := range ress {
			cam := rtspsimpleserver.SConf{}
			cam.Stream = stream
			cam.Conf = i.Conf

			res[stream] = cam
		}
	}

	if ctx.Err() != nil {
		return res, rtsp.err.SetError(ctx.Err())
	}

	return res, nil
}

// PostAddRTSP отправляет POST запрос на добавление потока
func (rtsp *repository) PostAddRTSP(ctx context.Context, cam rtspsimpleserver.SConf) ce.IError {

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
	fmt.Println("body", string(postJson))

	req, err := http.NewRequest(http.MethodPost, URLPostAdd, bytes.NewBuffer(postJson))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	req.Header.Add("Content-Type", "application/json")

	if ctx.Err() != nil {
		return rtsp.err.SetError(ctx.Err())
	}

	// Get запрос и обработка ошибки
	resp, err := rtsp.client.Do(req)
	if err != nil {
		return rtsp.err.SetError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rtsp.err.SetError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	return nil
}

// PostRemoveRTSP отправляет POST запрос на удаление потока
func (rtsp *repository) PostRemoveRTSP(ctx context.Context, camRTSP rtspsimpleserver.SConf) ce.IError {
	// Парсинг URL
	URLPostRemove := fmt.Sprintf(rtsp.cfg.Url + rtsp.cfg.Rtsp.Api.UrlRemove + camRTSP.Stream)
	rtsp.log.Debug("Url for request to rtsp:\n\t" + URLPostRemove)

	var buf []byte

	req, err := http.NewRequest(http.MethodPost, URLPostRemove, bytes.NewBuffer(buf))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	req.Header.Add("Content-Type", "application/json")

	if ctx.Err() != nil {
		return rtsp.err.SetError(ctx.Err())
	}

	// Get запрос и обработка ошибки
	resp, err := rtsp.client.Do(req)
	if err != nil {
		return rtsp.err.SetError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rtsp.err.SetError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	return nil
}

// PostEditRTSP отправляет POST запрос на изменение потока
func (rtsp *repository) PostEditRTSP(ctx context.Context, cam rtspsimpleserver.SConf) ce.IError {

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

	req, err := http.NewRequest(http.MethodPost, URLPostEdit, bytes.NewBuffer(postJson))
	if err != nil {
		return rtsp.err.SetError(err)
	}
	req.Header.Add("Content-Type", "application/json")

	if ctx.Err() != nil {
		return rtsp.err.SetError(ctx.Err())
	}

	// Get запрос и обработка ошибки
	resp, err := rtsp.client.Do(req)
	if err != nil {
		return rtsp.err.SetError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rtsp.err.SetError(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	return nil
}
