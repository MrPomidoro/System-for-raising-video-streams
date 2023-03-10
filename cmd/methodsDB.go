package service

import (
	"context"
	"errors"
	"fmt"

	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// insertIntoStatusStream принимает результат выполнения запроса через API (ошибка) и список камер с бд
// и выполняет вставку в таблицу status_stream
func (a *app) InsertIntoStatusStream(method string, ctx context.Context, cam rtsp.SConf, err ce.IError) ce.IError {
	if err != nil {
		a.log.Error(err.Error())

		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}

		insertStructStatusStream := statusstream.StatusStream{StreamId: cam.Id, StatusResponse: false}
		err = a.statusStreamRepo.Insert(ctx, &insertStructStatusStream)
		if err != nil {
			return a.err.SetError(fmt.Errorf("cannot insert stream %s to table status_stream", cam.Stream))
		}
		a.log.Debug("Success insert to table status_stream")

		return nil
	}

	a.log.Debug(fmt.Sprintf("Success complete post request for %s config %s", method, cam.Stream))

	if ctx.Err() != nil {
		return a.err.SetError(ctx.Err())
	}

	insertStructStatusStream := statusstream.StatusStream{StreamId: cam.Id, StatusResponse: true}
	err = a.statusStreamRepo.Insert(ctx, &insertStructStatusStream)
	if err == ce.ErrorStatusStream.SetError(errors.New("ERROR: insert or update on table \"status_stream\" violates foreign key constraint \"stream_id\" (SQLSTATE 23503)")) {
		return a.err.SetError(fmt.Errorf("cannot insert into status_stream: stream not in database"))
	}
	if err != nil {
		return err
	}

	a.log.Debug("Success insert to table status_stream")

	return nil
}
