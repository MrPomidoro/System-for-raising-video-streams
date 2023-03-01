package service

import (
	"context"
	"fmt"
	"os/signal"
	"reflect"
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
	defer close(a.doneChan)
	defer cancel()

	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sign := <-a.sigChan:
		a.log.Info(fmt.Sprintf("Got signal: %v, exiting", sign))
	case <-a.doneChan:
		a.log.Info("Found fatal error, exiting")
	}

	a.db.Close()
	a.log.Info("Close database connection")

	a.log.Debug("Waiting...")
}

// getDBAndApi реализует получение камер с базы данных и с rtsp
func (a *app) getDBAndApi(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream,
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
		runOnReady = fmt.Sprintf(cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)
	} else {
		runOnReady = ""
	}

	// Поле протокола не должно быть пустым
	// по умолчанию - tcp
	var protocol = camDB.Protocol.String
	if protocol == "" {
		protocol = "tcp"
	}

	return rtsp.SConf{
		Stream: camDB.Stream,
		Conf: rtsp.Conf{
			SourceProtocol: protocol,
			RunOnReady:     runOnReady,
			Source:         fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream),
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

// getCamsEdit - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая мапу камер, поля которых в бд и ртсп отличаются
func (a *app) getCamsEdit(dataDB []refreshstream.Stream, dataRTSP map[string]rtsp.SConf,
	camsAdd map[string]rtsp.SConf, camsRemove map[string]rtsp.SConf) map[string]rtsp.SConf {

	camsForEdit := make(map[string]rtsp.SConf)

	for _, camDB := range dataDB {

		cam := dbToCompare(a.cfg, camDB)
		// Проверяется, совпадают ли данные
		if reflect.DeepEqual(cam.Conf, rtspToCompare(dataRTSP[camDB.Stream])) {
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

// getCamsAdd - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая мапу камер, отсутствующих в rtsp, но имеющихся в базе
func (a *app) getCamsAdd(dataDB []refreshstream.Stream,
	dataRTSP map[string]rtsp.SConf) map[string]rtsp.SConf {

	camsForAdd := make(map[string]rtsp.SConf)

	for _, camDB := range dataDB {
		// Если камера есть в бд, но отсутствует в ртсп, добавляется в список
		if _, ok := dataRTSP[camDB.Stream]; ok {
			continue
		}
		cam := dbToCompare(a.cfg, camDB)
		camsForAdd[cam.Stream] = cam
	}

	return camsForAdd
}

// getCamsRemove - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// удаляющая из мапы с результатом из rtsp камеры, которые не нужно
func (a *app) getCamsRemove(dataDB []refreshstream.Stream,
	dataRTSP map[string]rtsp.SConf) {

	// fmt.Println("dataRTSP old", dataRTSP)
	for _, camDB := range dataDB {
		delete(dataRTSP, camDB.Stream)
	}
}
