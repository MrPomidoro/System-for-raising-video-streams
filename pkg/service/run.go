package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rsrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/repository"
	rsusecase "github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream/usecase"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ssrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository"
	ssusecase "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/usecase"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/rtsp"
	"github.com/sirupsen/logrus"
)

// Прототип приложения
type app struct {
	cfg                  *config.Config
	Log                  *logrus.Logger
	Db                   *sql.DB
	SigChan              chan os.Signal
	refreshStreamUseCase refreshstream.RefreshStreamUseCase
	statusStreamUseCase  statusstream.StatusStreamUseCase
}

// Функция, инициализирующая прототип приложения
func NewApp(cfg *config.Config) *app {
	log := logger.NewLog(cfg.LogLevel)
	db := database.CreateDBConnection(cfg)
	sigChan := make(chan os.Signal, 1)
	repoRS := rsrepository.NewRefreshStreamRepository(db)
	repoSS := ssrepository.NewStatusStreamRepository(db)

	return &app{
		cfg:                  cfg,
		Db:                   db,
		Log:                  log,
		SigChan:              sigChan,
		refreshStreamUseCase: rsusecase.NewRefreshStreamUseCase(repoRS, db, log),
		statusStreamUseCase:  ssusecase.NewStatusStreamUseCase(repoSS, db, log),
	}
}

// Алгоритм
func (a *app) Run() {

	// Инициализация контекста
	ctx := context.Background()

	go func() {
		// Канал для периодического выполнения алгоритма
		tick := time.NewTicker(time.Second * 7) //a.cfg.Refresh_Time)
		defer tick.Stop()
		for {
			fmt.Println("")

			// Выполняется периодически через установленный в конфигурационном файле промежуток времени
			<-tick.C

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, lenResDB, lenResRTSP, err := a.getDBAndApi(ctx)
			if err != nil {
				logger.LogError(a.Log, err)
				continue
			}

			// Сравнение числа записей в базе данных и записей в rtsp

			/*
				Если данных в базе столько же, сколько в rtsp:
				проверка, одинаковые ли записи:
				- если одинаковые, завершение и ожидание следующего запуска программы;
				- если различаются:
					- получение списка отличий,
					- отправка API,
					- запись в status_stream.
			*/
			if lenResDB == lenResRTSP {
				logger.LogInfo(a.Log, fmt.Sprintf("The count of data in the database = %d is equal to the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))

				// Проверка одинаковости данных по стримам
				identity := methods.CheckIdentity(dataDB, dataRTSP)

				if identity {
					logger.LogInfo(a.Log, "Data is identity, waiting...")
					continue
				}

				// Получение списков камер на добавление и удаление
				resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
				logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v",
					resSliceAdd, resSliceRemove))

				// Добавление камер
				a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
				// Удаление камер
				a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)

				//
				/*
					Если данных в базе больше, чем в rtsp:
					получение списка отличий;
					API на добавление в ртсп;
					запись в status_stream
				*/
			} else if lenResDB > lenResRTSP {
				logger.LogInfo(a.Log, fmt.Sprintf("The count of data in the database = %d is greater than the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))

				// Получение списков камер на добавление
				resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
				logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

				// Добавление камер
				a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
				// Удаление камер
				a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)

				//
				/*
					Если данных в базе меньше, чем в rtsp:
					получение списка отличий;
					API на добавление в ртсп;
					запись в status_stream
				*/
			} else if lenResDB < lenResRTSP {
				logger.LogInfo(a.Log, fmt.Sprintf("The count of data in the database = %d is less than the count of data in rtsp-simple-server = %d; waiting...", lenResDB, lenResRTSP))

				// Ожидание 5 секунд и повторный запрос данных с базы и с rtsp
				time.Sleep(time.Second * 5)
				_, _, lenResDBLESS, lenResRTSPLESS, err := a.getDBAndApi(ctx)
				if err != nil {
					logger.LogError(a.Log, err)
					continue
				}

				// Сравнение числа записей в базе данных и записей в rtsp после нового запроса
				if lenResDBLESS > lenResRTSPLESS {
					// Получение списков камер на добавление и удаление
					resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
					logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

					// Добавление камер
					a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
					// Удаление камер
					a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)

				} else if lenResDBLESS < lenResRTSPLESS {
					// Получение списков камер на добавление и удаление
					resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
					logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

					// Добавление камер
					a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
					// Удаление камер
					a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)

				} else if lenResDBLESS == lenResRTSPLESS {

					// Проверка одинаковости данных по стримам
					identity := methods.CheckIdentity(dataDB, dataRTSP)
					if identity {
						logger.LogInfo(a.Log, "Data is identity, waiting...")
						continue
					}

					// Получение списков камер на добавление и удаление
					resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
					logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v",
						resSliceAdd, resSliceRemove))

					// Добавление камер
					a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
					// Удаление камер
					a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
				}
			}
		}
	}()
}

/*
Метод для корректного завершения работы программы
при получении прерывающего сигнала
*/
func (a *app) GracefulShutdown(sig chan os.Signal) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	sign := <-sig
	logger.LogWarn(a.Log, fmt.Sprintf("Got signal: %v, exiting", sign))
	database.CloseDBConnection(a.cfg, a.Db)
	time.Sleep(time.Second * 2)
	close(a.SigChan)
}

// Get запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) []refreshstream.RefreshStream {
	req, err := a.refreshStreamUseCase.GetStatusTrue(ctx)
	if err != nil {
		logger.LogError(a.Log, err)
		return req
	}
	logger.LogDebug(a.Log, "Received response from the database")
	return req
}

/*
Получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, длины этих списков, статус код, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream,
	map[string]interface{}, int, int, error) {
	var lenResRTSP int

	// Отправка запросов к базе и к rtsp
	resDB := a.getReqFromDB(ctx)
	resRTSP := rtsp.GetRtsp(a.cfg)

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

			err := rtsp.PostAddRTSP(camDB, a.cfg)

			// Запись в базу данных результата выполнения
			if err != nil {
				logger.LogError(a.Log, err)
				insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
				err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
				if err != nil {
					logger.LogError(a.Log,
						"cannot insert to table status_stream")
					continue
				}
				logger.LogInfo(a.Log,
					"Success insert to table status_stream")
			}

			logger.LogInfo(a.Log, fmt.Sprintf("Success complete post request for add config %s", elemAdd))
			insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
			err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
			if err != nil {
				logger.LogError(a.Log,
					"cannot insert to table status_stream")
				continue
			}
			logger.LogInfo(a.Log,
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

	dataDB, err := a.refreshStreamUseCase.GetStatusFalse(ctx)
	if err != nil {
		logger.LogError(a.Log, err)
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

					err := rtsp.PostRemoveRTSP(camRTSP, a.cfg)

					// Запись в базу данных результата выполнения
					if err != nil {
						logger.LogError(a.Log, err)
						insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
						err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
						if err != nil {
							logger.LogError(a.Log,
								"cannot insert to table status_stream")
							continue
						}
						logger.LogInfo(a.Log,
							"Success insert to table status_stream")

					}

					logger.LogInfo(a.Log,
						fmt.Sprintf("Success complete Post request for remove config %s", elemRemove))
					insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
					err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
					if err != nil {
						logger.LogError(a.Log,
							"cannot insert to table status_stream")
						continue
					}
					logger.LogInfo(a.Log,
						"Success insert to table status_stream")

					break
				}
			}
		}
	}
}
