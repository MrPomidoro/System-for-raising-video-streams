package service

import (
	"context"
	"sync"
	"time"
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

	errChan := make(chan error)
	// Переподключение при необходимости
	go a.db.Ping(ctx, a.log, errChan)

	var mu sync.Mutex

loop:
	for {
		select {

		case <-ctx.Done():
			break loop

		// Выполняется периодически через установленный в конфигурационном файле промежуток времени
		case <-tick.C:
			if a.db.Db.Ping() != nil {
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

			// Получение отличающихся камер поля
			camsEdit := a.getCamsEdit(dataDB, dataRTSP)
			if len(camsEdit) == 0 {
				a.log.Info("Data is identity, waiting...")
				continue
			}

			// Если число камер совпадает, но стримы отличаются
			err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
			if err != nil {
				a.log.Error(err.Error())
				continue
			}

			// Если в бд и ртсп одни и те же камеры
			if isCamsSame(dataDB, dataRTSP) {
				// Если имеются отличия, отправляется запрос к ртсп на изменение
				a.log.Info("Count of data is same, but the values are different")
				err := a.editCamerasToRTSP(ctx, camsEdit)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}
			}
		}
	}
}
