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

/*
GracefulShutdown - метод для корректного завершения работы программы
при получении прерывающего сигнала
*/
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

/*
getDBAndApi реализует получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, ошибка
*/
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
			// resDB, err = a.getReqFromDB(ctx, dataDBchan)
			// a.getReqFromDB(ctx, dataDBchan)
			err = a.refreshStreamRepo.Get(ctx, true, dataDBchan)
			if err != nil {
				a.err.NextError(err)
				a.log.Error(a.err.Error())
				return
				// return nil, nil, a.err
			}
		}()

		// Отправка запроса к rtsp
		go func() {
			// resRTSP, err = a.rtspRepo.GetRtsp()
			a.rtspRepo.GetRtsp(ctx, dataRTSPchan)
			if err != nil {
				a.err.NextError(err)
				a.log.Error(a.err.Error())
				return
				// return nil, nil, a.err
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
				// fmt.Println("v from db") // тут всё ок
				resDB = append(resDB, v)
			}
		}
		fmt.Println("resDB", resDB)

	loop2:
		for {
			select {
			case <-ctx.Done():
				return resDB, resRTSP, nil
			case v, ok := <-dataRTSPchan:
				time.Sleep(100 * time.Millisecond)
				if !ok {
					fmt.Println("закрыто ртсп ваше")
					break loop2
				}
				fmt.Println("v from rtsp", v)
				mu.Lock()
				resRTSP[v.Stream] = v
				mu.Unlock()
			}
		}
	}

	return resDB, resRTSP, nil
}

/*
// equalOrIdentityData проверяет isEqualCount и identity:
// если оба - true, возвращает true;
// если isEqualCount - true, а identity - false, изменяет потоки и возвращает true;
// иначе возвращает false
func (a *app) equalOrIdentityData(ctx context.Context, isEqualCount, identity bool,
	sconfArr []rtspsimpleserver.SConf) bool {

	if isEqualCount && identity {
		a.log.Debug("Data is identity, waiting...")
		return true

	} else if isEqualCount && !identity {
		a.log.Debug("Count of data is same, but the field values are different")
		err := a.editCamerasToRTSP(ctx, sconfArr)
		if err != nil {
			a.log.Error(err.Error())
		}
		return true
	}
	return false
}
*/

// differentCount выполняется в случае, если число данных в базе и в rtsp отличается, возвращает ошибку при её наличии
func (a *app) differentCount(ctx context.Context, dataDB []refreshstream.RefreshStream, dataRTSP []rtspsimpleserver.SConf) ce.IError {
	a.log.Debug("Count of data is different")

	// err := a.addAndRemoveData(ctx, dataRTSP, dataDB)
	// if err != nil {
	// 	a.err.NextError(err)
	// 	return a.err
	// }
	return nil
}

//
//
//

// type SConfDB struct {
// 	Id int
// 	rtspsimpleserver.SConf
// }

func ConvertDBtoRTSP(cfg *config.Config, camDB refreshstream.RefreshStream) rtspsimpleserver.SConf {
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

// GetCamsAdd - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая мапу камер, отсутствующих в rtsp, но имеющихся в базе
func (a *app) getCamsAdd(dataDB []refreshstream.RefreshStream,
	dataRTSP map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {

	// camsForAdd := []rtspsimpleserver.SConf{}
	camsForAdd := make(map[string]rtspsimpleserver.SConf)

	for _, camDB := range dataDB {
		// Если камера есть в бд, но отсутствует в ртсп, добавляется в список
		if _, ok := dataRTSP[camDB.Stream.String]; ok {
			continue
		}
		cam := ConvertDBtoRTSP(a.cfg, camDB)
		camsForAdd[cam.Stream] = cam
		// camsForAdd = append(camsForAdd, cam)
	}

	return camsForAdd
}

// GetCamsEdit - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// возвращающая список камер, имеющихся в rtsp, но отсутствующих в базе
func (a *app) getCamsEdit(cfg *config.Config, dataDB []refreshstream.RefreshStream,
	dataRTSP map[string]rtspsimpleserver.SConf) map[string]rtspsimpleserver.SConf {

	camsForEdit := make(map[string]rtspsimpleserver.SConf)

	for _, camDB := range dataDB {

		cam := ConvertDBtoRTSP(cfg, camDB)
		// Проверяется, совпадают ли данные
		if reflect.DeepEqual(cam.Conf, dataRTSP[camDB.Stream.String]) {
			continue
		}
		// Если не совпадают, камера добавляется в мапу
		camsForEdit[cam.Stream] = cam
	}

	return camsForEdit
}

// GetCamsRemove - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
// удаляющая из мапы с результатом из rtsp камеры, которые не нужно
func getCamsRemove(dataDB []refreshstream.RefreshStream,
	dataRTSP map[string]rtspsimpleserver.SConf) {

	for _, camDB := range dataDB {
		delete(dataRTSP, camDB.Stream.String)
	}
}
