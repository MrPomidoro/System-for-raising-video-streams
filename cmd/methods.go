package service

import (
	"context"
	"fmt"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// GracefulShutdown - метод для корректного завершения работы программы
// при получении прерывающего сигнала
func (a *app) GracefulShutdown(cancel context.CancelFunc) {
	defer time.Sleep(time.Second * 5)
	defer close(a.sigChan)
	defer cancel()

	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)

	sign := <-a.sigChan
	a.log.Info(fmt.Sprintf("Got signal: %v, exiting", sign))

	a.db.Close()
	a.log.Info("Close database connection")

	a.log.Debug("Waiting...")
}

// getDBAndApi реализует получение камер с базы данных и с rtsp
func (a *app) GetDBAndApi(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream,
	map[string]rtsp.SConf, ce.IError) {

	// Отправка запроса к базе
	resDB, err := a.refreshStreamRepo.Get(ctx, true)
	if err != nil {
		return nil, nil, err
	}

	// Отправка запроса к rtsp
	resRTSP, err := a.rtspRepo.GetRtsp(ctx)
	if err != nil {
		return nil, nil, err
	}

	return resDB, resRTSP, nil
}

// dbToCompare приводит данные от бд к виду, который можно сравнить с ртсп
func dbToCompare(cfg *config.Config, camDB refreshstream.Stream) rtsp.SConf {
	// Парсинг поля RunOnReady
	var runOnReady string
	if cfg.Run != "" {
		runOnReady = fmt.Sprintf(cfg.Run, camDB.Port, camDB.CamPath.String)
	} else {
		runOnReady = ""
	}

	return rtsp.SConf{
		Stream: strings.TrimSpace(camDB.CodeMp),
		Conf: rtsp.Conf{
			SourceProtocol: "tcp",
			RunOnReady:     runOnReady,
			Source:         fmt.Sprintf("rtsp://%s:%s@%v:554/%s", strings.TrimSpace(camDB.Login.String), strings.TrimSpace(camDB.Pass.String), camDB.Ip.IPNet.IP, strings.TrimSpace(camDB.CamPath.String)),
		},
		Id: camDB.Id,
	}
}

// rtspToCompare приводит данные от ртсп к виду, который можно сравнить с бд
func rtspToCompare(camRTSP rtsp.SConf) rtsp.Conf {
	return rtsp.Conf{
		SourceProtocol: camRTSP.Conf.SourceProtocol,
		RunOnReady:     camRTSP.Conf.RunOnReady,
		Source:         camRTSP.Conf.Source,
	}
}

// GetCamsEdit - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая мапу камер, поля которых в бд и ртсп отличаются
func (a *app) GetCamsEdit(dataDB []refreshstream.Stream, dataRTSP map[string]rtsp.SConf,
	camsAdd map[string]rtsp.SConf, camsRemove map[string]rtsp.SConf) map[string]rtsp.SConf {

	camsForEdit := make(map[string]rtsp.SConf)

	for _, camDB := range dataDB {

		cam := dbToCompare(a.cfg, camDB)
		// Проверяется, совпадают ли данные
		if reflect.DeepEqual(cam.Conf, rtspToCompare(dataRTSP[camDB.CodeMp])) {
			continue
		}
		if _, ok := camsAdd[cam.Stream]; ok {
			continue
		}
		if _, ok := camsRemove[cam.Stream]; ok {
			continue
		}
		// Если камеры не совпадают и отсутствуют в списках на добавление и удаление,
		// камера добавляется в мапу
		camsForEdit[cam.Stream] = cam
	}

	return camsForEdit
}

// GetCamsAdd - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая мапу камер, отсутствующих в rtsp, но имеющихся в базе
func (a *app) GetCamsAdd(dataDB []refreshstream.Stream,
	dataRTSP map[string]rtsp.SConf) map[string]rtsp.SConf {

	camsForAdd := make(map[string]rtsp.SConf)

	for _, camDB := range dataDB {
		// Если камера есть в бд, но отсутствует в ртсп, добавляется в список
		if _, ok := dataRTSP[camDB.CodeMp]; ok {
			continue
		}
		cam := dbToCompare(a.cfg, camDB)
		camsForAdd[cam.Stream] = cam
	}

	return camsForAdd
}

// GetCamsRemove - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// удаляющая из мапы с результатом из rtsp камеры, которые не нужно
func (a *app) GetCamsRemove(dataDB []refreshstream.Stream,
	dataRTSP map[string]rtsp.SConf) {

	for _, camDB := range dataDB {
		delete(dataRTSP, camDB.CodeMp)
	}
}
