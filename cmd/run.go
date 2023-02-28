package service

import (
	"context"
	"sync"
	"time"

	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
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

	var mu sync.Mutex

loop:
	for {
		select {

		case <-ctx.Done():
			break loop

		// Выполняется периодически через установленный в конфигурационном файле промежуток времени
		case <-tick.C:
			if a.db.Conn.Ping(ctx) != nil {
				continue loop
			}

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, err := a.getDBAndApi(ctx, &mu)
			if err != nil {
				a.log.Error(err.Error())
				continue
			}

			if ctx.Err() != nil {
				continue
			}

			camsRemove := make(map[string]rtspsimpleserver.SConf)
			transcode.Transcode(dataRTSP, &camsRemove)
			a.getCamsRemove(dataDB, camsRemove)

			camsAdd := a.getCamsAdd(dataDB, dataRTSP)

			camsEdit := a.getCamsEdit(dataDB, dataRTSP, camsAdd, camsRemove)

			if len(camsEdit) == 0 && len(camsRemove) == 0 && len(camsAdd) == 0 {
				a.log.Info("Data is identity, waiting...")
				continue
			}

			// Если число камер совпадает, но стримы отличаются
			// err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
			err = a.addAndRemoveData(ctx, dataDB, dataRTSP, camsAdd, camsRemove)
			if err != nil {
				a.log.Error(err.Error())
				continue
			}

			// Если в бд и ртсп одни и те же камеры
			if len(camsEdit) != 0 {
				// Если имеются отличия, отправляется запрос к ртсп на изменение
				a.log.Info("Cameras is same, but its values are different")
				err := a.editCamerasToRTSP(ctx, camsEdit)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}
			}
		}
	}
}
