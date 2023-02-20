package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
)

/*
addCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо добавить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на добавление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) addCamerasToRTSP(ctx context.Context, resSliceAdd []string,
	dataDB []refreshstream.RefreshStream) ce.IError {
	// Перебор всех элементов списка камер на добавление
	for _, elemAdd := range resSliceAdd {
		// Цикл для извлечения данных из структуры выбранной камеры
		for _, camDB := range dataDB {
			if camDB.Stream.String != elemAdd {
				continue
			}

			err := a.rtspRepo.PostAddRTSP(camDB)
			if err != nil {
				a.err.NextError(err.GetError())
				return a.err
			}

			err = a.refreshStreamRepo.Update(ctx, camDB.Stream.String)
			if err != nil {
				a.err.NextError(err.GetError())
				return a.err
			}
			a.log.Debug("Success send request to update stream_status")

			// Запись в базу данных результата выполнения
			err = a.insertIntoStatusStream("add", ctx, camDB, err)
			if err != nil {
				a.err.NextError(err.GetError())
				return a.err
			}
		}
	}
	return nil
}

/*
removeCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо удалить
с rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на удаление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) removeCamerasToRTSP(ctx context.Context, resSliceRemove []string,
	dataRTSP []rtspsimpleserver.SConf) ce.IError {

	dataDB, err := a.refreshStreamRepo.Get(ctx, false)
	if err != nil {
		a.err.NextError(err.GetError())
		return a.err
	}

	// Цикл для извлечения данных из структуры выбранной камеры
	for _, camRTSP := range dataRTSP {
		// Для возможности извлечения данных
		// camsRTSPMap := camsRTSP.(map[string]interface{})

		// camRTSP - стрим камеры
		// for camRTSP := range camsRTSPMap {

		// Перебор всех камер, которые нужно удалить
		for _, elemRemove := range resSliceRemove {
			if camRTSP.Stream != elemRemove {
				continue
			}

			for _, camDB := range dataDB {
				if camDB.Stream.String != elemRemove {
					continue
				}

				err := a.rtspRepo.PostRemoveRTSP(camRTSP.Stream)
				if err != nil {
					a.err.NextError(err.GetError())
					return a.err
				}

				// Запись в базу данных результата выполнения
				err = a.insertIntoStatusStream("remove", ctx, camDB, err)
				if err != nil {
					a.err.NextError(err.GetError())
					return a.err
				}
			}
		}
	}
	// }
	return nil
}

/*
editCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо изменить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на изменение камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) editCamerasToRTSP(ctx context.Context, confArr []rtspsimpleserver.SConf,
	dataDB []refreshstream.RefreshStream) ce.IError {
	for _, camDB := range dataDB {
		for _, sconf := range confArr {

			if camDB.Stream.String != sconf.Stream {
				continue
			}

			if sconf.Conf.SourceProtocol == "" && sconf.Conf.Source == "" && (sconf.Conf.RunOnReady == "" && a.cfg.Run != "") {
				continue
			}

			err := a.rtspRepo.PostEditRTSP(camDB, sconf)
			if err != nil {
				a.err.NextError(err.GetError())
				return a.err
			}

			// Запись в базу данных результата выполнения
			err = a.insertIntoStatusStream("edit", ctx, camDB, err)
			if err != nil {
				a.err.NextError(err.GetError())
				return a.err
			}
		}
	}
	return nil
}

/*
addAndRemoveData - метод, в которым выполняются функции, получающие списки
отличающихся данных, выполняется удаление лишних камер и добавление недостающих
*/
func (a *app) addAndRemoveData(ctx context.Context, dataRTSP []rtspsimpleserver.SConf,
	dataDB []refreshstream.RefreshStream) ce.IError {

	// Получение списков камер на добавление и удаление
	resSliceAdd := methods.GetCamsForAdd(dataDB, dataRTSP)
	resSliceRemove := methods.GetCamsForRemove(dataDB, dataRTSP)

	a.log.Debug(fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

	// Добавление камер
	if resSliceAdd != nil {
		err := a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
		if err != nil {
			a.err.NextError(err.GetError())
			return a.err
		}
	}

	// Удаление камер
	if resSliceRemove != nil {
		err := a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
		if err != nil {
			a.err.NextError(err.GetError())
			return a.err
		}
	}
	return nil
}
