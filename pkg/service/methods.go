package service

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
)

/*
GracefulShutdown - метод для корректного завершения работы программы
при получении прерывающего сигнала
*/
func (a *app) GracefulShutdown(ctx context.Context, cancel context.CancelFunc) {

	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)
	sign := <-a.sigChan

	a.log.Warn(fmt.Sprintf("Got signal: %v, exiting", sign))
	cancel()
	database.CloseDBConnection(a.cfg, a.Db)
	a.log.Error(ctx.Err().Error())
	time.Sleep(time.Second * 10)
	close(a.sigChan)
}

/*
Используется в API
insertIntoStatusStream принимает результат выполнения запроса через API (ошибка) и список камер с бд
и выполняет вставку в таблицу status_stream
*/
func (a *app) insertIntoStatusStream(method string, ctx context.Context, camDB refreshstream.RefreshStream, err error) error {
	if err != nil {
		a.log.Error(err.Error())
		insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
		err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
		if err != nil {
			a.log.Error("cannot insert to table status_stream")
			return err
		}
		a.log.Info("Success insert to table status_stream")
	}

	a.log.Info(fmt.Sprintf("Success complete post request for %s config %s", method, camDB.Stream.String))
	insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
	err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
	if err != nil {
		a.log.Error("cannot insert to table status_stream")
		return err
	}
	a.log.Info("Success insert to table status_stream")

	return nil
}

/*
addAndRemoveData - метод, в которым выполняются функции, получающие списки
отличающихся данных, выполняется удаление лишних камер и добавление недостающих
*/
func (a *app) addAndRemoveData(ctx context.Context, dataRTSP map[string]interface{},
	dataDB []refreshstream.RefreshStream) error {

	// Получение списков камер на добавление и удаление
	resSliceAdd := methods.GetCamsForAdd(dataDB, dataRTSP)
	resSliceRemove := methods.GetCamsForRemove(dataDB, dataRTSP)

	a.log.Info(fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

	// Добавление камер
	if resSliceAdd != nil {
		err := a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
		if err != nil {
			return err
		}
	}

	// Удаление камер
	if resSliceRemove != nil {
		err := a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
equalOrIdentityData проверяет isEqualCount и identity:
если оба - true, возвращает true;
если isEqualCount - true, а identity - false, изменяет потоки и возвращает true;
иначе возвращает false
*/
func (a *app) equalOrIdentityData(ctx context.Context, isEqualCount, identity bool,
	confArr []rtspsimpleserver.Conf, dataDB []refreshstream.RefreshStream) bool {

	if isEqualCount && identity {
		a.log.Info("Data is identity, waiting...")
		return true

	} else if isEqualCount && !identity {
		a.log.Info("Count of data is same, but the field values are different")
		err := a.editCamerasToRTSP(ctx, confArr, dataDB)
		if err != nil {
			a.log.Error(err.Error())
		}
		return true
	}
	return false
}

// differentCount выполняется в случае, если число данных в базе и в rtsp, возвращает ошибку при её наличии
func (a *app) differentCount(ctx context.Context, dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) error {
	a.log.Error("Count of data is different")
	err := a.addAndRemoveData(ctx, dataRTSP, dataDB)
	if err != nil {
		return err
	}
	return nil
}
