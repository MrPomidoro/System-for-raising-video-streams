package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// getReqFromDB реализует Get запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) ([]refreshstream.RefreshStream, *ce.Error) {
	req, err := a.refreshStreamRepo.Get(ctx, true)
	if err != nil {
		a.err.NextError(err)
		return req, a.err
	}
	if len(req) == 0 {
		return nil, a.err.SetError(fmt.Errorf("response from database not received"))
	}

	a.log.Debug("Received response from the database")
	return req, nil
}

/*
Используется в API
insertIntoStatusStream принимает результат выполнения запроса через API (ошибка) и список камер с бд
и выполняет вставку в таблицу status_stream
*/
func (a *app) insertIntoStatusStream(method string, ctx context.Context, camDB refreshstream.RefreshStream, err *ce.Error) *ce.Error {
	if err != nil {
		a.log.Error(err.Error())
		insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
		err = a.statusStreamRepo.Insert(ctx, &insertStructStatusStream)
		if err != nil {
			a.log.Error("cannot insert to table status_stream")
			a.err.NextError(err)
			return a.err
		}
		a.log.Debug("Success insert to table status_stream")

		return nil
	}

	a.log.Debug(fmt.Sprintf("Success complete post request for %s config %s", method, camDB.Stream.String))
	insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
	err = a.statusStreamRepo.Insert(ctx, &insertStructStatusStream)
	if err != nil {
		a.log.Error("cannot insert to table status_stream")
		a.err.NextError(err)
		return a.err
	}
	a.log.Debug("Success insert to table status_stream")

	return nil
}
