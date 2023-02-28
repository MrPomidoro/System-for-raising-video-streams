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

	a.db.CloseDBConnection()

	a.log.Debug("Waiting...")
}

// getDBAndApi реализует получение камер с базы данных и с rtsp
func (a *app) getDBAndApi(ctx context.Context, mu *sync.Mutex) ([]refreshstream.Stream,
	map[string]rtspsimpleserver.SConf, ce.IError) {

	var err ce.IError

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
	// return resDB, resRTSP, a.err.SetError(errors.New("ашипка"))
}

// isCamsSame проверяет, что камеры в бд и в ртсп совпадают, т.е. учитывает
// случай, когда количество камер равно, но сами камеры отличаются;
// напр., камеры в бд: 1, 2, в ртсп: 1, 3 --- камеру 2 добавить, камеру 3 удалить.
// Возвращает true, если камеры совпадают, false - если отличаются
func isCamsSame(dataDB []refreshstream.Stream, dataRTSP map[string]rtspsimpleserver.SConf) bool {
	counter := 0
	for _, camDB := range dataDB {
		for camRTSP := range dataRTSP {
			if camDB.Stream.String == camRTSP {
				counter++
			}
		}
	}
	fmt.Println(counter, len(dataDB))
	return counter == len(dataDB)
}

// dbToCompare приводит данные от бд к виду, который можно сравнить с ртсп
func dbToCompare(cfg *config.Config, camDB refreshstream.Stream) rtspsimpleserver.SConf {
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
// возвращающая мапу камер, поля которых в бд и ртсп отличаются
func (a *app) getCamsEdit(dataDB []refreshstream.Stream,
	dataRTSP map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {

	camsForEdit := make(map[string]rtspsimpleserver.SConf)

	for _, camDB := range dataDB {

		cam := dbToCompare(a.cfg, camDB)
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
func (a *app) getCamsAdd(dataDB []refreshstream.Stream,
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
func (a *app) getCamsRemove(dataDB []refreshstream.Stream,
	dataRTSP map[string]rtspsimpleserver.SConf) {

	for _, camDB := range dataDB {
		delete(dataRTSP, camDB.Stream.String)
	}
}
