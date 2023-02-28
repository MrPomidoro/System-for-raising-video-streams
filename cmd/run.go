package service

import (
	"context"
	"fmt"
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
			// if a.db.Ping() != nil {
			// 	continue loop
			// }

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, err := a.getDBAndApi(ctx, &mu)
			if err != nil {
				a.log.Error(err.Error())
				continue
			}

			if ctx.Err() != nil {
				continue
			}

			// Получение отличающихся камер поля
			camsEdit := a.getCamsEdit(dataDB, dataRTSP)
			fmt.Println("camsForEdit", camsEdit)

			camsRemove := make(map[string]rtspsimpleserver.SConf)
			transcode.CopyMap(dataRTSP)
			fmt.Println("camsRemove addAndRemoveData", camsRemove)
			if len(camsEdit) == 0 && len(camsRemove) == 0 {
				a.log.Info("Data is identity, waiting...")
				continue
			}

			// Если число камер совпадает, но стримы отличаются
			// err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
			err = a.addAndRemoveData(ctx, camsRemove, dataDB)
			if err != nil {
				a.log.Error(err.Error())
				continue
			}

			// Если в бд и ртсп одни и те же камеры
			if isCamsSame(dataDB, dataRTSP) && len(camsEdit) != 0 {
				// Если имеются отличия, отправляется запрос к ртсп на изменение
				a.log.Info("Count of data is same, but the values are different")
				fmt.Println("camsForEdit", camsEdit)
				err := a.editCamerasToRTSP(ctx, camsEdit)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}
			}
		}
	}
}
