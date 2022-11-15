package service

import (
	"context"
	"database/sql"
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
		tick := time.NewTicker(time.Second * 3) //a.cfg.Refresh_Time)
		defer tick.Stop()
		for {
			fmt.Println("")
			select {
			// Выполняется периодически через установленный в конфигурационном файле промежуток времени
			case <-tick.C:

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
					identity := checkIdentity(dataDB, dataRTSP)

					if identity {
						logger.LogInfo(a.Log, "Data is identity, waiting...")
						continue
					} else {
						resSliceAdd, resSliceRemove := getDifferenceElements(dataDB, dataRTSP)
						logger.LogInfo(a.Log, fmt.Sprintf("Elements to be added: %v --- Elements to be removed: %v", resSliceAdd, resSliceRemove))

						continue
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

					resSliceAdd, resSliceRemove := getDifferenceElements(dataDB, dataRTSP)
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
		}
	}()
}

// Метод для корректного завершения работы программы
// при получении прерывающего сигнала
func (a *app) GracefulShutdown(sig chan os.Signal) {
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sign := <-sig:
		logger.LogWarn(a.Log, fmt.Sprintf("Got signal: %v, exiting", sign))
		time.Sleep(time.Second * 2)
		database.CloseDBConnection(a.cfg, a.Db)
		close(a.SigChan)
	}
}
