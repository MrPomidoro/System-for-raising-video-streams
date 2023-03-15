package service

import (
	"context"
	"sync"
	"time"

	rtsp "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/transcode"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~ //
//   ~~~   Алгоритм   ~~~   //
// ~~~~~~~~~~~~~~~~~~~~~~~~ //

func (a *app) Run(ctx context.Context) {
	a.log.Info("Start service")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Канал для периодического выполнения алгоритма
	tick := time.NewTicker(a.cfg.RefreshTime)
	defer tick.Stop()

	// Проверка коннекта к базе данных
	if a.db == nil {
		return
	}

	// Создаем канал для получения оповещений о сбое подключения
	errCh := make(chan error)

	// Запускаем асинхронную проверку поддержания соединения
	go a.db.KeepAlive(ctx, a.log, errCh)

	var mu sync.Mutex

loop:
	for {

		select {

		case <-ctx.Done():
			break loop

		// Выполняется периодически через установленный в конфигурационном файле промежуток времени
		case <-tick.C:

			if !a.db.IsConn(ctx) {
				continue loop
			}

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, err := a.GetDBAndApi(ctx, &mu)
			if err != nil {
				a.log.Error(err.Error())
				continue loop
			}

			if ctx.Err() != nil {
				break loop
			}

			camsRemove := make(map[string]rtsp.SConf)
			transcode.Transcode(dataRTSP, &camsRemove)
			a.GetCamsRemove(dataDB, camsRemove)

			camsAdd := a.GetCamsAdd(dataDB, dataRTSP)

			camsEdit := a.GetCamsEdit(dataDB, dataRTSP, camsAdd, camsRemove)

			if len(camsEdit) == 0 && len(camsRemove) == 0 && len(camsAdd) == 0 {
				a.log.Info("Data is identity, waiting...")
				continue loop
			}

			err = a.AddRemoveData(ctx, dataDB, dataRTSP, camsAdd, camsRemove)
			if err != nil {
				a.log.Error(err.Error())
				continue loop
			}

			err = a.EditData(ctx, camsEdit)
			if err != nil {
				a.log.Error(err.Error())
			}
		}
	}
}
