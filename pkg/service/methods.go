package service

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
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
		a.log.Info("Found ERROR, exiting")
	}

	a.db.CloseDBConnection(a.cfg)

	a.log.Debug("sleep...")
}

/*
getDBAndApi реализует получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream,
	[]rtspsimpleserver.SConf, ce.IError) {
	var resRTSP []rtspsimpleserver.SConf
	var resDB []refreshstream.RefreshStream

	// dataDB := make(chan refreshstream.RefreshStream)
	// dataRTSP := make(chan)

	// Отправка запроса к базе
	resDB, err := a.getReqFromDB(ctx)
	if err != nil {
		a.err.NextError(err)
		return nil, nil, a.err
	}

	// Отправка запроса к rtsp
	resRTSP, err = a.rtspRepo.GetRtsp()
	if err != nil {
		a.err.NextError(err)
		return nil, nil, a.err
	}

	return resDB, resRTSP, nil
}

/*
equalOrIdentityData проверяет isEqualCount и identity:
если оба - true, возвращает true;
если isEqualCount - true, а identity - false, изменяет потоки и возвращает true;
иначе возвращает false
*/
func (a *app) equalOrIdentityData(ctx context.Context, isEqualCount, identity bool,
	sconfArr []rtspsimpleserver.SConf, dataDB []refreshstream.RefreshStream) bool {

	if isEqualCount && identity {
		a.log.Debug("Data is identity, waiting...")
		return true

	} else if isEqualCount && !identity {
		a.log.Debug("Count of data is same, but the field values are different")
		err := a.editCamerasToRTSP(ctx, sconfArr, dataDB)
		if err != nil {
			a.log.Error(err.Error())
		}
		return true
	}
	return false
}

// differentCount выполняется в случае, если число данных в базе и в rtsp отличается, возвращает ошибку при её наличии
func (a *app) differentCount(ctx context.Context, dataDB []refreshstream.RefreshStream, dataRTSP []rtspsimpleserver.SConf) ce.IError {
	a.log.Debug("Count of data is different")

	err := a.addAndRemoveData(ctx, dataRTSP, dataDB)
	if err != nil {
		a.err.NextError(err)
		return a.err
	}
	return nil
}
