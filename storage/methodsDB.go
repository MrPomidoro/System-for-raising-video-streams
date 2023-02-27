package service

import (
	"context"
	"fmt"

	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	statusstream "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// getReqFromDB реализует Get запрос на получение списка камер из базы данных,
// данные отправляются в канал
/*
func (a *app) getReqFromDB(ctx context.Context, dataDBchan chan refreshstream.RefreshStream) ce.IError {
	err := a.refreshStreamRepo.Get(ctx, true, dataDBchan)
	if err != nil {
		a.err.NextError(err)
		return a.err
	}
	// if len(req) == 0 {
	// return  a.err.SetError(fmt.Errorf("no response from database received"))
	// return nil, a.err.SetError(fmt.Errorf("no response from database received"))
	// }

	a.log.Debug("Received response from the database")
	return nil
}
*/

/*
Используется в API
insertIntoStatusStream принимает результат выполнения запроса через API (ошибка) и список камер с бд
и выполняет вставку в таблицу status_stream
*/
func (a *app) insertIntoStatusStream(method string, ctx context.Context, cam rtspsimpleserver.SConf, err ce.IError) ce.IError {
	if err != nil {
		a.log.Error(err.Error())
		insertStructStatusStream := statusstream.StatusStream{StreamId: cam.Id, StatusResponse: false}
		err = a.statusStreamRepo.Insert(ctx, &insertStructStatusStream)
		if err != nil {
			a.err.NextError(err)
			return a.err.SetError(fmt.Errorf("cannot insert stream %s to table status_stream", cam.Stream))
		}
		a.log.Debug("Success insert to table status_stream")

		return nil
	}

	a.log.Debug(fmt.Sprintf("Success complete post request for %s config %s", method, cam.Stream))
	insertStructStatusStream := statusstream.StatusStream{StreamId: cam.Id, StatusResponse: true}
	err = a.statusStreamRepo.Insert(ctx, &insertStructStatusStream)
	if err != nil {
		a.err.NextError(err)
		return a.err.SetError(fmt.Errorf("cannot insert to table status_stream"))
	}
	a.log.Debug("Success insert to table status_stream")

	return nil
}
