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
	repoRS := rsrepository.NewRefreshStreamRepository(db, log)
	repoSS := ssrepository.NewStatusStreamRepository(db, log)

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
func (a *app) Run() error {
	ctx := context.Background()
	logger.LogDebug(a.Log, "Context initializated")

	// Число потоков после выполнения запроса к rtsp
	var lenResRTSP int
	// Отправка запросов к базе и к rtsp
	resDB := a.getReqFromDB(ctx)
	resRTSP := rtsp.GetRtsp(a.cfg)

	a.getReqFromRtsp()

	// Определение числа потоков с rtsp
	for _, items := range resRTSP { // items - поле "items"
		// Для доступа к данным каждой камеры:
		camsMap := items.(map[string]interface{})
		lenResRTSP = len(camsMap)
	}

	if len(resDB) == lenResRTSP {
		logger.LogInfo(a.Log, fmt.Sprintf("The number of cameras in the data = %d is equal to the number of data in RTSP = %d\n", len(resDB), lenResRTSP))
		EqualData()
	} else if len(resDB) > lenResRTSP {
		logger.LogInfo(a.Log, fmt.Sprintf("The number of cameras in the data = %d is less than the number of data in RTSP = %d\n", len(resDB), lenResRTSP))
		LessData()
	} else if len(resDB) < lenResRTSP {
		logger.LogInfo(a.Log, fmt.Sprintf("The number of cameras in the data = %d is greater than the number of data in RTSP = %d\n", len(resDB), lenResRTSP))
		MoreData()
	}

	//
	// ssExample := statusstream.StatusStream{StreamId: 3, StatusResponse: true}
	// // Запись в базу данных результата выполнения (нужно менять)
	// err = a.statusStreamUseCase.Insert(ctx, &ssExample)
	// if err != nil {
	// 	logger.LogError(a.Log, "cannot insert")
	// } else {
	// 	logger.LogDebug(a.Log, "insert correct, 200")
	// }

	return nil
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
