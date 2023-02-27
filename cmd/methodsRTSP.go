package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
)

/*
addCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо добавить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на добавление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) addCamerasToRTSP(ctx context.Context, camsAdd map[string]rtspsimpleserver.SConf) ce.IError {

	// Перебор всех элементов списка камер на добавление
	for _, camAdd := range camsAdd {
		// Цикл для извлечения данных из структуры выбранной камеры
		// for _, camDB := range dataDB {
		// if camDB.Stream.String != elemAdd {
		// 	continue
		// }

		err := a.rtspRepo.PostAddRTSP(camAdd)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}

		err = a.refreshStreamRepo.Update(ctx, camAdd.Stream)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
		a.log.Debug("Success send request to update stream_status")

		// Запись в базу данных результата выполнения
		err = a.insertIntoStatusStream("add", ctx, camAdd, err)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
	}
	return nil
}

/*
removeCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо удалить
с rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на удаление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) removeCamerasToRTSP(ctx context.Context, dataRTSP map[string]rtspsimpleserver.SConf) ce.IError {

	// Перебор всех камер, которые нужно удалить
	for _, cam := range dataRTSP {
		// Цикл для извлечения данных из структуры выбранной камеры
		// for _, camDB := range dataDB {
		// if camDB.Stream.String != elemRemove {
		// 	continue
		// }

		err := a.rtspRepo.PostRemoveRTSP(cam)
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
		// }
	}
	// }
	return nil
}

/*
editCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо изменить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на изменение камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) editCamerasToRTSP(ctx context.Context, camsForEdit map[string]rtspsimpleserver.SConf) ce.IError {

	for _, cam := range camsForEdit {

		if cam.Conf.SourceProtocol == "" && cam.Conf.Source == "" && (cam.Conf.RunOnReady == "" && a.cfg.Run != "") {
			continue
		}

		err := a.rtspRepo.PostEditRTSP(cam)
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

/*
addAndRemoveData - метод, в которым выполняются функции, получающие списки
отличающихся данных, выполняется удаление лишних камер и добавление недостающих
*/
func (a *app) addAndRemoveData(ctx context.Context, dataRTSP map[string]rtspsimpleserver.SConf,
	dataDB []refreshstream.RefreshStream) ce.IError {

	// Получение списков камер на добавление и удаление
	camsAdd := a.getCamsAdd(dataDB, dataRTSP)
	getCamsRemove(dataDB, dataRTSP)

	a.log.Debug(fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", camsAdd, dataRTSP))

	// Добавление камер
	if camsAdd != nil {
		err := a.addCamerasToRTSP(ctx, camsAdd)
		if err != nil {
			a.err.NextError(err)
			return a.err
		}
	}

	// Удаление камер
	// if resSliceRemove != nil {
	err := a.removeCamerasToRTSP(ctx, dataRTSP)
	if err != nil {
		a.err.NextError(err)
		return a.err
	}
	// }
	return nil
}
