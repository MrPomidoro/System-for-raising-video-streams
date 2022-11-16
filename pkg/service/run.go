package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
				if resSliceAdd != nil {
					a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
				}
				// Удаление камер
				if resSliceRemove != nil {
					a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
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

				// Получение списков камер на добавление
				resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
				logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

				// Добавление камер
				if resSliceAdd != nil {
					a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
				}
				// Удаление камер
				if resSliceRemove != nil {
					a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
				}

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
				dataDB, dataRTSP, lenResDBLESS, lenResRTSPLESS, err := a.getDBAndApi(ctx)
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
					if resSliceAdd != nil {
						a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
					}
					// Удаление камер
					if resSliceRemove != nil {
						a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
					}

				} else if lenResDBLESS < lenResRTSPLESS {
					// Получение списков камер на добавление и удаление
					resSliceAdd, resSliceRemove := methods.GetDifferenceElements(dataDB, dataRTSP)
					logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

					// Добавление камер
					if resSliceAdd != nil {
						a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
					}
					// Удаление камер
					if resSliceRemove != nil {
						a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
					}

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
					if resSliceAdd != nil {
						a.addCamerasToRTSP(ctx, resSliceAdd, dataDB)
					}
					// Удаление камер
					if resSliceRemove != nil {
						a.removeCamerasToRTSP(ctx, resSliceRemove, dataRTSP)
					}
				}
			}
		}
	}()
}
