package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
)

// getReqFromDB реализует Get запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) ([]refreshstream.RefreshStream, error) {
	req, err := a.refreshStreamRepo.Get(ctx, true)
	if err != nil {
		return req, err
	}
	a.log.Debug("Received response from the database")

	return req, nil
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
		a.log.Debug("Success insert to table status_stream")

		return nil
	}

	a.log.Debug(fmt.Sprintf("Success complete post request for %s config %s", method, camDB.Stream.String))
	insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
	err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
	if err != nil {
		a.log.Error("cannot insert to table status_stream")
		return err
	}
	a.log.Debug("Success insert to table status_stream")

	return nil
}
