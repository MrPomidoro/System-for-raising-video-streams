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
	LogStatusCode        *logrus.Logger
	Db                   *sql.DB
	SigChan              chan os.Signal
	refreshStreamUseCase refreshstream.RefreshStreamUseCase
	statusStreamUseCase  statusstream.StatusStreamUseCase
}

// Функция, инициализирующая прототип приложения
func NewApp(cfg *config.Config) *app {
	log := logger.NewLog(cfg.LogLevel)
	logStatCode := logger.NewLogStatCode(cfg.LogLevel)
	db := database.CreateDBConnection(cfg)
	sigChan := make(chan os.Signal, 1)
	repoRS := rsrepository.NewRefreshStreamRepository(db, logStatCode)
	repoSS := ssrepository.NewStatusStreamRepository(db)

	return &app{
		cfg:                  cfg,
		Db:                   db,
		Log:                  log,
		LogStatusCode:        logStatCode,
		SigChan:              sigChan,
		refreshStreamUseCase: rsusecase.NewRefreshStreamUseCase(repoRS, db, log),
		statusStreamUseCase:  ssusecase.NewStatusStreamUseCase(repoSS, db, log),
	}
}

// Алгоритм
func (a *app) Run() {
	// Инициализация контекста
	ctx := context.Background()
	logger.LogDebug(a.Log, "Context initializated")

	go func() {
		// Канал для периодического выполнения алгоритма
		tick := time.NewTicker(time.Second * 5) //a.cfg.Refresh_Time)
		defer tick.Stop()
		for {
			fmt.Println("")

			// Выполняется периодически через установленный в конфигурационном файле промежуток времени
			<-tick.C

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, lenResDB, lenResRTSP, stCode, err := a.getDBAndApi(ctx)
			if err != nil {
				logger.LogErrorStatusCode(a.LogStatusCode, err, "Get", stCode)
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
				resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
				logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v",
					resSliceAdd, resSliceRemove))

				// Перебор всех камер, которые нужно добавить
				for _, elemAdd := range resSliceAdd {
					// Цикл для извлечения данных из структуры выбранной камеры
					for _, camDB := range dataDB {
						if camDB.Stream.String != elemAdd {
							continue
						}

						rtsp.PostAddRTSP(camDB, a.cfg) // err =
						fmt.Println("add", elemAdd)

						// /*
						// Запись в базу данных результата выполнения
						if err != nil {
							logger.LogErrorStatusCode(a.LogStatusCode, fmt.Sprintf("cannot complete post request: %v", err), "Post", "500")
							insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
							err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
							if err != nil {
								logger.LogErrorStatusCode(a.LogStatusCode,
									"cannot insert to table status_stream", "Post", "400")
								continue
							}
							logger.LogInfoStatusCode(a.LogStatusCode,
								"Success insert to table status_stream", "Post", "200")

							continue
						}
						// Запись в базу данных результата выполнения
						insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
						err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
						if err != nil {
							logger.LogErrorStatusCode(a.LogStatusCode,
								"cannot insert to table status_stream", "Post", "400")
							continue
						}
						logger.LogInfoStatusCode(a.LogStatusCode,
							"Success insert to table status_stream", "Post", "200")

						// */
					}
				}

				// Перебор всех камер, которые нужно удалить
				for _, elemRemove := range resSliceRemove {
					// Цикл для извлечения данных из структуры выбранной камеры
					for _, camDB := range dataDB {
						if camDB.Stream.String == elemRemove {
							continue
						}

						rtsp.PostRemoveRTSP(camDB, a.cfg) // err =
						fmt.Println("remove", elemRemove)

						// /*
						// Запись в базу данных результата выполнения
						if err != nil {
							logger.LogErrorStatusCode(a.LogStatusCode, fmt.Sprintf("cannot complete post request: %v", err), "Post", "500")
							insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: false}
							err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
							if err != nil {
								logger.LogErrorStatusCode(a.LogStatusCode,
									"cannot insert to table status_stream", "Post", "400")
								continue
							}
							logger.LogInfoStatusCode(a.LogStatusCode,
								"Success insert to table status_stream", "Post", "200")

							continue
						}
						// Запись в базу данных результата выполнения
						insertStructStatusStream := statusstream.StatusStream{StreamId: camDB.Id, StatusResponse: true}
						err = a.statusStreamUseCase.Insert(ctx, &insertStructStatusStream)
						if err != nil {
							logger.LogErrorStatusCode(a.LogStatusCode,
								"cannot insert to table status_stream", "Post", "400")
							continue
						}
						logger.LogInfoStatusCode(a.LogStatusCode,
							"Success insert to table status_stream", "Post", "200")

						// */

						break
					}
				}

				//
				/*
					Если данных в базе больше, чем в rtsp:
					получение списка отличий;
					API на добавление в ртсп;
					запись в status_stream
				*/
			} else if lenResDB > lenResRTSP {
				logger.LogInfo(a.Log, fmt.Sprintf("The count of data in the database = %d is greater than the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))

				resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
				logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

				//
				/*
					Если данных в базе меньше, чем в rtsp:
					получение списка отличий;
					API на добавление в ртсп;
					запись в status_stream
				*/
			} else if lenResDB < lenResRTSP {
				logger.LogInfo(a.Log, fmt.Sprintf("The count of data in the database = %d is less than the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))

				// Ожидание 5 секунд и повторный запрос данных с базы и с rtsp
				time.Sleep(time.Second * 5)
				_, _, lenResDBLESS, lenResRTSPLESS, stCode, err := a.getDBAndApi(ctx)
				if err != nil {
					logger.LogErrorStatusCode(a.LogStatusCode, err, "Get", stCode)
					continue
				}

				// Сравнение числа записей в базе данных и записей в rtsp после нового запроса
				if lenResDBLESS > lenResRTSPLESS {
					/*
						получаем список отличий;
						апи на добавление в ртсп;
						запись в статус_стрим
					*/
					continue
				} else if lenResDBLESS < lenResRTSPLESS {
					/*
						получаем список отличий;
						апи на удаление в ртсп;
						запись в статус_стрим
					*/
					continue
				}
			}

			/*
				ssExample := statusstream.StatusStream{StreamId: 3, StatusResponse: true}
				// Запись в базу данных результата выполнения (нужно менять)
				err = a.statusStreamUseCase.Insert(ctx, &ssExample)
				if err != nil {
					logger.LogErrorStatusCode(a.LogStatusCode, "cannot insert", "Post", "400")
				} else {
					logger.LogInfoStatusCode(a.LogStatusCode, "Success insert", "Post", "200")
				}
			*/

		}
	}()
}

// Метод для корректного завершения работы программы
// при получении прерывающего сигнала
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
	req, err := a.refreshStreamUseCase.Get(ctx)
	if err != nil {
		logger.LogError(a.Log, fmt.Sprintf("cannot get response from database: %v", err))
	}
	return req
}

/*
Получение списка камер с базы данных и с rtsp
На выходе: список с бд, список с rtsp, длины этих списков, статус код, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream, map[string]interface{}, int, int, string, error) {
	var lenResRTSP int

	// Отправка запросов к базе и к rtsp
	resDB := a.getReqFromDB(ctx)
	resRTSP := rtsp.GetRtsp(a.cfg)

	// Проверка, что ответ от базы данных не пустой
	if len(resDB) == 0 {
		return resDB, resRTSP, len(resDB), lenResRTSP, "400", errors.New("response from database is null")
	}

	// Определение числа потоков с rtsp
	for _, items := range resRTSP { // items - поле "items"
		// мапа: ключ - номер камеры, значения - остальные поля этой камеры
		camsMap := items.(map[string]interface{})
		lenResRTSP = len(camsMap) // количество камер
	}

	// Проверка, что ответ от rtsp данных не пустой
	if lenResRTSP == 0 {
		return resDB, resRTSP, len(resDB), lenResRTSP, "500", errors.New("response from rtsp-simple-server is null")
	}

	return resDB, resRTSP, len(resDB), lenResRTSP, "200", nil
}
