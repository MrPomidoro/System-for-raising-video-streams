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
	var err ce.IError

	dataDBchan := make(chan refreshstream.RefreshStream)
	dataRTSPchan := make(chan rtspsimpleserver.SConf)

	select {
	case <-ctx.Done():
	default:
		// Отправка запроса к базе
		go func() {
			// resDB, err = a.getReqFromDB(ctx, dataDBchan)
			a.getReqFromDB(ctx, dataDBchan)
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

	loop:
		for {
			select {
			case <-ctx.Done():
				return resDB, resRTSP, nil

			case v, ok := <-dataDBchan:
				if !ok {
					break loop
				}
				fmt.Println("v from db") // тут всё ок
				resDB = append(resDB, v)
			}
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
			resRTSP = append(resRTSP, v)
		}
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
