package service

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

// addAndRemoveData - метод, в которым выполняются функции, получающие списки
// отличающихся данных, выполняется удаление лишних камер и добавление недостающих
func (a *app) AddRemoveData(ctx context.Context, dataDB []refreshstream.Stream,
	dataRTSP map[string]rtsp.SConf, camsAdd map[string]rtsp.SConf,
	camsRemove map[string]rtsp.SConf) ce.IError {

	if len(camsAdd) != 0 || len(camsRemove) != 0 {
		a.log.Info("The cameras are different")
	} else {
		return nil
	}

	// Добавление камер
	if len(camsAdd) != 0 {
		err := a.AddData(ctx, camsAdd)
		if err != nil {
			a.log.Error(err.Error())
		}
	}

	// Удаление камер
	if len(camsRemove) != 0 {
		err := a.RemoveData(ctx, camsRemove)
		if err != nil {
			a.log.Error(err.Error())
		}
	}

	return nil
}

// addCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо добавить
// в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на добавление камер,
// добавляет в таблицу status_stream запись с результатом выполнения запроса
func (a *app) AddData(ctx context.Context, camsAdd map[string]rtsp.SConf) ce.IError {

	// Перебор всех элементов списка камер на добавление
	for _, camAdd := range camsAdd {
		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}

		err := a.rtspRepo.PostAddRTSP(ctx, camAdd)
		if err != nil {
			return err
		}

		// err = a.refreshStreamRepo.Update(ctx, camAdd.Stream)
		// if err != nil {
		// 	return err
		// }
		// a.log.Debug("Success send request to update stream_status")

		// Запись в базу данных результата выполнения
		err = a.InsertIntoStatusStream("add", ctx, camAdd, err)
		if err != nil {
			return err
		}
	}
	return nil
}

// removeCamerasFromRTSP - функция, принимающая на вход список камер, которые необходимо удалить
// с rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на удаление камер,
// добавляет в таблицу status_stream запись с результатом выполнения запроса
func (a *app) RemoveData(ctx context.Context, dataRTSP map[string]rtsp.SConf) ce.IError {

	// Перебор всех камер, которые нужно удалить
	for _, cam := range dataRTSP {

		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}

		err := a.rtspRepo.PostRemoveRTSP(ctx, cam)
		if err != nil {
			return err
		}

		// Запись в базу данных результата выполнения
		err = a.InsertIntoStatusStream("remove", ctx, cam, err)
		if err != nil {
			return err
		}
	}

	return nil
}

// editData - функция, принимающая на вход список камер, которые необходимо изменить
// в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на изменение камер,
// добавляет в таблицу status_stream запись с результатом выполнения запроса
func (a *app) EditData(ctx context.Context, camsEdit map[string]rtsp.SConf) ce.IError {

	// Если в бд и ртсп одни и те же камеры
	if len(camsEdit) == 0 {
		return nil
	}

	// Если имеются отличия, отправляется запрос к ртсп на изменение
	a.log.Info("Cameras is same, but its values are different")

	for _, cam := range camsEdit {

		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}

		err := a.rtspRepo.PostEditRTSP(ctx, cam)
		if err != nil {
			return err
		}

		// Запись в базу данных результата выполнения
		err = a.InsertIntoStatusStream("edit", ctx, cam, err)
		if err != nil {
			return err
		}
	}

	return nil
}
