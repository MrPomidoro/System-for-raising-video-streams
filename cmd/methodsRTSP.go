package service

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
)

/*
addAndRemoveData - метод, в которым выполняются функции, получающие списки
отличающихся данных, выполняется удаление лишних камер и добавление недостающих
*/
func (a *app) addAndRemoveData(ctx context.Context, dataRTSP map[string]rtspsimpleserver.SConf,
	dataDB []refreshstream.Stream) ce.IError {

	// Получение мапы камер на добавление
	camsAdd := a.getCamsAdd(dataDB, dataRTSP)
	// Получение мапы камер на удаление; создаётся копия исходной мапы
	// с помощью метода Transcode, чтобы исходная не изменялась
	dataRTSPCopy := make(map[string]rtspsimpleserver.SConf)
	methods.Transcode(dataRTSP, dataRTSPCopy)
	a.getCamsRemove(dataDB, dataRTSPCopy)

	// Добавление камер
	if camsAdd != nil {
		err := a.addCamerasToRTSP(ctx, camsAdd)
		if err != nil {
			return a.err
		}
	}

	// Удаление камер
	if len(dataRTSPCopy) != 0 {
		err := a.removeCamerasFromRTSP(ctx, dataRTSPCopy)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
	}

	return nil
}

/*
addCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо добавить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на добавление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) addCamerasToRTSP(ctx context.Context, camsAdd map[string]rtspsimpleserver.SConf) ce.IError {

	// Перебор всех элементов списка камер на добавление
	for _, camAdd := range camsAdd {
		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}
		err := a.rtspRepo.PostAddRTSP(ctx, camAdd)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}

		// err = a.refreshStreamRepo.Update(ctx, camAdd.Stream)
		// if err != nil {
		// 	a.err.NextError(err)
		// 	return a.err
		// }
		// a.log.Debug("Success send request to update stream_status")

		// Запись в базу данных результата выполнения
		err = a.insertIntoStatusStream("add", ctx, camAdd, err)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
	}
	return nil
}

// removeCamerasFromRTSP - функция, принимающая на вход список камер, которые необходимо удалить
// с rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на удаление камер,
// добавляет в таблицу status_stream запись с результатом выполнения запроса
func (a *app) removeCamerasFromRTSP(ctx context.Context, dataRTSP map[string]rtspsimpleserver.SConf) ce.IError {

	// Перебор всех камер, которые нужно удалить
	for _, cam := range dataRTSP {

		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}

		err := a.rtspRepo.PostRemoveRTSP(ctx, cam)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}

		// Запись в базу данных результата выполнения
		err = a.insertIntoStatusStream("remove", ctx, cam, err)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
	}

	return nil
}

// editCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо изменить
// в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на изменение камер,
// добавляет в таблицу status_stream запись с результатом выполнения запроса
func (a *app) editCamerasToRTSP(ctx context.Context, camsForEdit map[string]rtspsimpleserver.SConf) ce.IError {

	for _, cam := range camsForEdit {

		// if cam.Conf.SourceProtocol == "" && cam.Conf.Source == "" && (cam.Conf.RunOnReady == "" && a.cfg.Run != "") {
		// 	continue
		// }

		if ctx.Err() != nil {
			return a.err.SetError(ctx.Err())
		}

		err := a.rtspRepo.PostEditRTSP(ctx, cam)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}

		// Запись в базу данных результата выполнения
		err = a.insertIntoStatusStream("edit", ctx, cam, err)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
	}

	return nil
}
