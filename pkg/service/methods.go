package service

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
)

/*
GracefulShutdown - метод для корректного завершения работы программы
при получении прерывающего сигнала
*/
func (a *app) GracefulShutdown(ctx context.Context, cancel context.CancelFunc) {

	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)
	sign := <-a.sigChan

	a.log.Info(fmt.Sprintf("Got signal: %v, exiting", sign))
	cancel()

	database.DBI.CloseDBConnection(a.db, a.cfg)
	a.log.Debug(ctx.Err().Error())

	a.log.Debug("sleep...")
	time.Sleep(time.Second * 10)
	close(a.sigChan)
}

/*
getDBAndApi реализует получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream,
	map[string]interface{}, error) {
	var resRTSP map[string]interface{}
	var resDB []refreshstream.RefreshStream

	// Отправка запроса к базе
	resDB, err := a.getReqFromDB(ctx)
	if err != nil {
		return []refreshstream.RefreshStream{}, map[string]interface{}{}, err
	}
	a.log.Debug("Get response from database")

	// Отправка запроса к rtsp
	resRTSP, err = a.rtspUseCase.GetRtsp()
	if err != nil {
		return []refreshstream.RefreshStream{}, map[string]interface{}{}, err
	}
	a.log.Debug("Get response from rtsp-simple-server")

	return resDB, resRTSP, nil
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
		a.log.Debug("Data is identity, waiting...")
		return true

	} else if isEqualCount && !identity {
		a.log.Debug("Count of data is same, but the field values are different")
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
	a.log.Debug("Count of data is different")
	err := a.addAndRemoveData(ctx, dataRTSP, dataDB)
	if err != nil {
		return err
	}
	return nil
}
