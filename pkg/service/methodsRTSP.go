package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
)

/*
addCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо добавить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на добавление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) addCamerasToRTSP(ctx context.Context, resSliceAdd []string,
	dataDB []refreshstream.RefreshStream) error {
	// Перебор всех элементов списка камер на добавление
	for _, elemAdd := range resSliceAdd {
		// Цикл для извлечения данных из структуры выбранной камеры
		for _, camDB := range dataDB {
			if camDB.Stream.String != elemAdd {
				continue
			}

			err := a.rtspUseCase.PostAddRTSP(camDB)
			if err != nil {
				return err
			}

			err = a.refreshStreamRepo.Update(ctx, camDB.Stream.String)
			if err != nil {
				return err
			}
			a.log.Debug("Success send request to update stream_status")

			// Запись в базу данных результата выполнения
			err = a.insertIntoStatusStream("add", ctx, camDB, err)
			if err != nil {
				return err
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
	dataRTSP map[string]interface{}) error {

	dataDB, err := a.refreshStreamRepo.Get(ctx, false)
	if err != nil {
		return err
	}

	// Цикл для извлечения данных из структуры выбранной камеры
	for _, camsRTSP := range dataRTSP {
		// Для возможности извлечения данных
		camsRTSPMap := camsRTSP.(map[string]interface{})

		// camRTSP - стрим камеры
		for camRTSP := range camsRTSPMap {

			// Перебор всех камер, которые нужно удалить
			for _, elemRemove := range resSliceRemove {
				if camRTSP != elemRemove {
					continue
				}

				for _, camDB := range dataDB {
					if camDB.Stream.String != elemRemove {
						continue
					}

					err := a.rtspUseCase.PostRemoveRTSP(camRTSP)
					if err != nil {
						return err
					}

					// Запись в базу данных результата выполнения
					err = a.insertIntoStatusStream("remove", ctx, camDB, err)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

/*
editCamerasToRTSP - функция, принимающая на вход список камер, которые необходимо изменить
в rtsp-simple-server, и список камер из базы данных. Отправляет Post запрос к rtsp на изменение камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) editCamerasToRTSP(ctx context.Context, confArr []rtspsimpleserver.Conf,
	dataDB []refreshstream.RefreshStream) error {
	for _, camDB := range dataDB {
		for _, conf := range confArr {

			if camDB.Stream.String != conf.Stream {
				continue
			}

			if conf.SourceProtocol == "" && conf.Source == "" && (conf.RunOnReady == "" && a.cfg.Run != "") {
				continue
			}

			err := a.rtspUseCase.PostEditRTSP(camDB, conf)
			if err != nil {
				return err
			}

			// Запись в базу данных результата выполнения
			err = a.insertIntoStatusStream("edit", ctx, camDB, err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
addAndRemoveData - метод, в которым выполняются функции, получающие списки
отличающихся данных, выполняется удаление лишних камер и добавление недостающих
*/
func (a *app) addAndRemoveData(ctx context.Context, dataRTSP map[string]interface{},
	dataDB []refreshstream.RefreshStream) error {

	// Получение списков камер на добавление и удаление
	resSliceAdd := methods.GetCamsForAdd(dataDB, dataRTSP)
	resSliceRemove := methods.GetCamsForRemove(dataDB, dataRTSP)

	a.log.Debug(fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

	// Добавление камер
	if resSliceAdd != nil {
		err := a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
		if err != nil {
			return err
		}
	}

	// Удаление камер
	if resSliceRemove != nil {
		err := a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
		if err != nil {
			return err
		}
	}
	return nil
}
