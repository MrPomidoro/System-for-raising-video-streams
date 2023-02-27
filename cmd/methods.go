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
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
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

	a.db.CloseDBConnection(a.cfg)

	a.log.Debug("Waiting...")
}

// getDBAndApi реализует получение списка камер с базы данных и с rtsp
// На выходе: список с бд, мапа с rtsp, ошибка
func (a *app) getDBAndApi(ctx context.Context, mu *sync.Mutex) ([]refreshstream.RefreshStream,
	map[string]rtspsimpleserver.SConf, ce.IError) {

	resRTSP := make(map[string]rtspsimpleserver.SConf)
	var resDB []refreshstream.RefreshStream
	var err ce.IError

	dataDBchan := make(chan refreshstream.RefreshStream)
	dataRTSPchan := make(chan rtspsimpleserver.SConf)

	select {
	case <-ctx.Done():
	default:
		// Отправка запроса к базе
		go func() {
			err = a.refreshStreamRepo.Get(ctx, true, dataDBchan)
			if err != nil {
				a.err.NextError(err)
				a.log.Error(a.err.Error())
				return
			}
		}()

		// Отправка запроса к rtsp
		go func() {
			a.rtspRepo.GetRtsp(ctx, dataRTSPchan)
			if err != nil {
				a.err.NextError(err)
				a.log.Error(a.err.Error())
				return
			}
		}()

		// Чтение
	loop:
		for {
			select {
			case <-ctx.Done():
				return resDB, resRTSP, nil

			case v, ok := <-dataDBchan:
				if !ok {
					break loop
				}
				resDB = append(resDB, v)
			}
		}

	loop2:
		for {
			select {
			case <-ctx.Done():
				return resDB, resRTSP, nil
			case v, ok := <-dataRTSPchan:
				time.Sleep(100 * time.Millisecond)
				if !ok {
					break loop2
				}
				mu.Lock()
				resRTSP[v.Stream] = v
				mu.Unlock()
			}
		}
	}

	return resDB, resRTSP, nil
}

//
//
//

// dbToCompare приводит данные от бд к виду, который можно сравнить с ртсп
func dbToCompare(cfg *config.Config, camDB refreshstream.RefreshStream) rtspsimpleserver.SConf {
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

	return rtspsimpleserver.SConf{
		Stream: camDB.Stream.String,
		Conf: rtspsimpleserver.Conf{
			SourceProtocol: protocol,
			RunOnReady:     runOnReady,
			Source:         fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream.String),
		},
		Id: camDB.Id,
	}
}

// rtspToCompare приводит данные от ртсп к виду, который можно сравнить с бд
func rtspToCompare(camRTSP rtspsimpleserver.SConf) rtspsimpleserver.Conf {
	return rtspsimpleserver.Conf{
		SourceProtocol: camRTSP.Conf.SourceProtocol,
		RunOnReady:     camRTSP.Conf.RunOnReady,
		Source:         camRTSP.Conf.Source,
	}
}

// GetCamsEdit - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая список камер, имеющихся в rtsp, но отсутствующих в базе
func (a *app) getCamsEdit(cfg *config.Config, dataDB []refreshstream.RefreshStream,
	dataRTSP map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {

	camsForEdit := make(map[string]rtspsimpleserver.SConf)

	for _, camDB := range dataDB {

		cam := dbToCompare(cfg, camDB)
		// Проверяется, совпадают ли данные
		if reflect.DeepEqual(cam.Conf, rtspToCompare(dataRTSP[camDB.Stream.String])) {
			continue
		}
		// Если не совпадают, камера добавляется в мапу
		camsForEdit[cam.Stream] = cam
	}

	return camsForEdit
}

// GetCamsAdd - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая мапу камер, отсутствующих в rtsp, но имеющихся в базе
func (a *app) getCamsAdd(dataDB []refreshstream.RefreshStream,
	dataRTSP map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {

	camsForAdd := make(map[string]rtspsimpleserver.SConf)

	for _, camDB := range dataDB {
		// Если камера есть в бд, но отсутствует в ртсп, добавляется в список
		if _, ok := dataRTSP[camDB.Stream.String]; ok {
			continue
		}
		cam := dbToCompare(a.cfg, camDB)
		camsForAdd[cam.Stream] = cam
	}

	return camsForAdd
}

// GetCamsRemove - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// удаляющая из мапы с результатом из rtsp камеры, которые не нужно
func getCamsRemove(dataDB []refreshstream.RefreshStream,
	dataRTSP map[string]rtspsimpleserver.SConf) {

	for _, camDB := range dataDB {
		delete(dataRTSP, camDB.Stream.String)
	}
}
