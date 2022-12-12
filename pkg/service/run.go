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
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	rtsprepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/repository"
	rtspusecase "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server/usecase"
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream"
	ssrepository "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/repository"
	ssusecase "github.com/Kseniya-cha/System-for-raising-video-streams/internal/statusstream/usecase"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
	"github.com/sirupsen/logrus"
)

// app - прототип приложения
type app struct {
	cfg                  *config.Config
	log                  *logrus.Logger
	Db                   *sql.DB
	SigChan              chan os.Signal
	refreshStreamUseCase refreshstream.RefreshStreamUseCase
	statusStreamUseCase  statusstream.StatusStreamUseCase
	rtspUseCase          rtspsimpleserver.RTSPUseCase
}

// NewApp инициализирует прототип приложения
func NewApp(cfg *config.Config) *app {
	log := logger.NewLog(cfg.LogLevel)
	if !cfg.Database_Connect {
		logger.LogError(log, "no permission to connect to database")
		return &app{}
	}
	db := database.CreateDBConnection(cfg)
	sigChan := make(chan os.Signal, 1)
	repoRS := rsrepository.NewRefreshStreamRepository(db)
	repoSS := ssrepository.NewStatusStreamRepository(db)
	repoRTSP := rtsprepository.NewRTSPRepository(cfg, log)

	return &app{
		cfg:                  cfg,
		Db:                   db,
		log:                  log,
		SigChan:              sigChan,
		refreshStreamUseCase: rsusecase.NewRefreshStreamUseCase(repoRS, db),
		statusStreamUseCase:  ssusecase.NewStatusStreamUseCase(repoSS, db),
		rtspUseCase:          rtspusecase.NewRTSPUseCase(repoRTSP),
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~ //
//   ~~~   Алгоритм   ~~~   //
// ~~~~~~~~~~~~~~~~~~~~~~~~ //

func (a *app) Run(ctx context.Context) {
	ctx, _ = context.WithCancel(ctx)

	// Канал для периодического выполнения алгоритма
	tick := time.NewTicker(a.cfg.Refresh_Time)
	defer tick.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-tick.C:
			// Выполняется периодически через установленный в конфигурационном файле промежуток времени
			fmt.Println("")

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, err := a.getDBAndApi(ctx)
			if err != nil {
				logger.LogError(a.log, err)
				continue
			}
			lenResDB, lenResRTSP := methods.CheckEmptyData(dataDB, dataRTSP)

			// ---------------------------------------------------------- //
			//   Сравнение числа записей в базе данных и записей в rtsp   //
			// ---------------------------------------------------------- //

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
				logger.LogInfo(a.log, fmt.Sprintf("The count of data in the database = %d is equal to the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))

				// Проверка одинаковости данных по стримам
				isEqualCount, identity, confArr := methods.CheckIdentityAndCountOfData(dataDB, dataRTSP, a.cfg)

				// Если число данных совпадает и данные одинаковые ИЛИ если число данных совпадает, но данные отличаются,
				// метод equalOrIdentityData возвращает true
				eqId := a.equalOrIdentityData(ctx, isEqualCount, identity, confArr, dataDB)
				if eqId {
					continue
				}

				// Если число данных отличается, выполняется differentCount
				err := a.differentCount(ctx, dataDB, dataRTSP)
				if err != nil {
					logger.LogError(a.log, err)
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

				logger.LogInfo(a.log, fmt.Sprintf("The count of data in the database = %d is greater than the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))
				err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
				if err != nil {
					logger.LogError(a.log, err)
					continue
				}

				//
				/*
					Если данных в базе меньше, чем в rtsp:
					получение списка отличий;
					API на добавление в ртсп;
					запись в status_stream
				*/
			} else if lenResDB < lenResRTSP {
				logger.LogInfo(a.log, fmt.Sprintf("The count of data in the database = %d is less than the count of data in rtsp-simple-server = %d; waiting...", lenResDB, lenResRTSP))

				// Ожидание 5 секунд и повторный запрос данных с базы и с rtsp
				time.Sleep(time.Second * 5)
				dataDB, dataRTSP, err := a.getDBAndApi(ctx)
				if err != nil {
					logger.LogError(a.log, err)
					continue
				}
				lenResDBLESS, lenResRTSPLESS := methods.CheckEmptyData(dataDB, dataRTSP)

				// Сравнение числа записей в базе данных и записей в rtsp после нового запроса
				if lenResDBLESS == lenResRTSPLESS {

					// Проверка одинаковости данных по стримам
					isEqualCount, identity, confArr := methods.CheckIdentityAndCountOfData(dataDB, dataRTSP, a.cfg)

					// Если число данных совпадает и данные одинаковые ИЛИ если число данных совпадает, но данные отличаются,
					// метод equalOrIdentityData возвращает true
					eqId := a.equalOrIdentityData(ctx, isEqualCount, identity, confArr, dataDB)
					if eqId {
						continue
					}

					// Если число данных отличается, выполняется differentCount
					err := a.differentCount(ctx, dataDB, dataRTSP)
					if err != nil {
						logger.LogError(a.log, err)
						continue
					}

				} else {

					err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
					if err != nil {
						logger.LogError(a.log, err)
						continue
					}

				}
			}
		}
	}
}
