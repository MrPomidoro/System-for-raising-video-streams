package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

/*
Метод для корректного завершения работы программы
при получении прерывающего сигнала
*/
func (a *app) GracefulShutdown(sig chan os.Signal) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	sign := <-sig
	logger.LogWarn(a.log, fmt.Sprintf("Got signal: %v, exiting", sign))
	database.CloseDBConnection(a.cfg, a.Db)
	time.Sleep(time.Second * 2)
	close(a.SigChan)
}

// Get запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) []refreshstream.RefreshStream {
	req, err := a.refreshStreamUseCase.Get(ctx, true)
	if err != nil {
		logger.LogError(a.log, err)
		return req
	}
	logger.LogDebug(a.log, "Received response from the database")
	return req
}

/*
Получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, длины этих списков, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream,
	map[string]interface{}, int, int, error) {
	var lenResRTSP int

	// Отправка запросов к базе и к rtsp
	resDB := a.getReqFromDB(ctx)
	resRTSP := a.rtspUseCase.GetRtsp()

	// Проверка, что ответ от базы данных не пустой
	if len(resDB) == 0 {
		return resDB, resRTSP, len(resDB), lenResRTSP, errors.New("response from database is null")
	}

	// Определение числа потоков с rtsp
	for _, items := range resRTSP { // items - поле "items"
		// мапа: ключ - номер камеры, значения - остальные поля этой камеры
		camsMap := items.(map[string]interface{})
		lenResRTSP = len(camsMap) // количество камер
	}

	// Проверка, что ответ от rtsp данных не пустой
	if lenResRTSP == 0 {
		return resDB, resRTSP, len(resDB), lenResRTSP, errors.New("response from rtsp-simple-server is null")
	}

	return resDB, resRTSP, len(resDB), lenResRTSP, nil
}

/*
Функция, принимающая на вход список камер, которые необходимо добавить в rtsp-simple-server,
и список камер из базы данных. Отправляет Post запрос к rtsp на добавление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) addCamerasToRTSP(ctx context.Context, resSliceAdd []string,
	dataDB []refreshstream.RefreshStream) {
	// Перебор всех элементов списка камер на добавление
	for _, elemAdd := range resSliceAdd {
		// Цикл для извлечения данных из структуры выбранной камеры
		for _, camDB := range dataDB {
			if camDB.Stream.String != elemAdd {
				continue
			}

			err := a.rtspUseCase.PostAddRTSP(camDB)

			// Запись в базу данных результата выполнения

			if err != nil {
				logger.LogError(a.log, err)
				insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
				err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
				if err != nil {
					logger.LogError(a.log,
						"cannot insert to table status_stream")
					continue
				}
				logger.LogInfo(a.log,
					"Success insert to table status_stream")
			}

			logger.LogInfo(a.log, fmt.Sprintf("Success complete post request for add config %s", elemAdd))

			insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
			err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
			if err != nil {
				logger.LogError(a.log,
					"cannot insert to table status_stream")
				continue
			}
			logger.LogInfo(a.log,
				"Success insert to table status_stream")
		}
	}
}

/*
Функция, принимающая на вход список камер, которые необходимо удалить из rtsp-simple-server,
и список камер из базы данных. Отправляет Post запрос к rtsp на удаление камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) removeCamerasToRTSP(ctx context.Context, resSliceRemove []string,
	dataRTSP map[string]interface{}) {

	dataDB, err := a.refreshStreamUseCase.Get(ctx, false)
	if err != nil {
		logger.LogError(a.log, err)
		return
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

					// Запись в базу данных результата выполнения

					if err != nil {
						logger.LogError(a.log, err)
						insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
						err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
						if err != nil {
							logger.LogError(a.log,
								"cannot insert to table status_stream")
							continue
						}
						logger.LogInfo(a.log,
							"Success insert to table status_stream")

					}

					logger.LogInfo(a.log,
						fmt.Sprintf("Success complete Post request for remove config %s", elemRemove))
					insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
					err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
					if err != nil {
						logger.LogError(a.log,
							"cannot insert to table status_stream")
						continue
					}
					logger.LogInfo(a.log,
						"Success insert to table status_stream")

					break
				}
			}
		}
	}
}

/*
Функция, принимающая на вход список камер, которые необходимо изменить в rtsp-simple-server,
и список камер из базы данных. Отправляет Post запрос к rtsp на изменение камер,
добавляет в таблицу status_stream запись с результатом выполнения запроса
*/
func (a *app) editCamerasToRTSP(ctx context.Context, confArr []rtspsimpleserver.Conf,
	dataDB []refreshstream.RefreshStream) {
	for _, camDB := range dataDB {
		for _, conf := range confArr {
			if camDB.Stream.String != conf.Stream {
				continue
			}

			if (conf.SourceProtocol == camDB.Protocol.String || conf.SourceProtocol == "") && conf.RunOnReady == "" {
				continue
			}

			err := a.rtspUseCase.PostEditRTSP(camDB, conf)

			// Запись в базу данных результата выполнения

			if err != nil {
				logger.LogError(a.log, err)
				insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
				err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
				if err != nil {
					logger.LogError(a.log,
						"cannot insert to table status_stream")
					continue
				}
				logger.LogInfo(a.log,
					"Success insert to table status_stream")
			}

			logger.LogInfo(a.log, fmt.Sprintf("Success complete post request for edit config %s", camDB.Stream.String))
			insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
			err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
			if err != nil {
				logger.LogError(a.log,
					"cannot insert to table status_stream")
				continue
			}
			logger.LogInfo(a.log,
				"Success insert to table status_stream")
		}
	}
}
