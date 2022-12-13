package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
)

/*
GracefulShutdown - метод для корректного завершения работы программы
при получении прерывающего сигнала
*/
func (a *app) GracefulShutdown(sig chan os.Signal, ctx context.Context, cancel context.CancelFunc) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	sign := <-sig
	a.log.Warn(fmt.Sprintf("Got signal: %v, exiting", sign))
	cancel()
	database.CloseDBConnection(a.cfg, a.Db)
	a.log.Error(ctx.Err().Error())
	time.Sleep(time.Second * 10)
	close(a.SigChan)
}

// ------------------------------------------------- //
//                    Get запросы                    //
// ------------------------------------------------- //

// getReqFromDB реализует Get запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) ([]refreshstream.RefreshStream, error) {
	req, err := a.refreshStreamUseCase.Get(ctx, true)
	if err != nil {
		return req, err
	}
	a.log.Debug("Received response from the database")
	return req, nil
}

/*
getDBAndApi реализует получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream,
	map[string]interface{}, error) {
	var resRTSP map[string]interface{}
	var resDB []refreshstream.RefreshStream

	// Отправка запроса к базе
	resDB, err := a.getReqFromDB(ctx)
	if err != nil {
		return []refreshstream.RefreshStream{}, map[string]interface{}{}, err
	}

	// Отправка запроса к rtsp
	resRTSP, err = a.rtspUseCase.GetRtsp()
	if err != nil {
		return []refreshstream.RefreshStream{}, map[string]interface{}{}, err
	}

	return resDB, resRTSP, nil
}

// ----------------------------------------- //
//                    API                    //
// ---------------------------------------- //

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

			err = a.refreshStreamUseCase.Update(ctx, camDB.Stream.String)
			if err != nil {
				return err
			} else {
				a.log.Debug("Success send request to update stream_status")
			}

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

	dataDB, err := a.refreshStreamUseCase.Get(ctx, false)
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

// --------------------------------------------------- //
//                    Другие методы                    //
// --------------------------------------------------- //

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
		a.log.Info("Success insert to table status_stream")
	}

	a.log.Info(fmt.Sprintf("Success complete post request for %s config %s", method, camDB.Stream.String))
	insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
	err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
	if err != nil {
		a.log.Error("cannot insert to table status_stream")
		return err
	}
	a.log.Info("Success insert to table status_stream")

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

	a.log.Info(fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

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

//
//
//

/*
equalOrIdentityData проверяет isEqualCount и identity:
если оба - true, возвращает true;
если isEqualCount - true, а identity - false, изменяет потоки и возвращает true;
иначе возвращает false
*/
func (a *app) equalOrIdentityData(ctx context.Context, isEqualCount, identity bool,
	confArr []rtspsimpleserver.Conf, dataDB []refreshstream.RefreshStream) bool {

	if isEqualCount && identity {
		a.log.Info("Data is identity, waiting...")
		return true

	} else if isEqualCount && !identity {
		a.log.Info("Count of data is same, but the field values are different")
		err := a.editCamerasToRTSP(ctx, confArr, dataDB)
		if err != nil {
			a.log.Error(err.Error())
		}
		return true
	}
	return false
}

// differentCount выполняется в случае, если число данных в базе и в rtsp, возвращает ошибку при её наличии
func (a *app) differentCount(ctx context.Context, dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) error {
	a.log.Error("Count of data is different")
	err := a.addAndRemoveData(ctx, dataRTSP, dataDB)
	if err != nil {
		return err
	}
	return nil
}
