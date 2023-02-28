package service

import (
	"context"
	"fmt"
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

			a.log.Info(fmt.Sprintf("The count of data in the database = %d is equal to the count of data in rtsp-simple-server = %d", len(dataDB), len(dataRTSP)))

			// Получение отличающихся камер поля
			camsForEdit := a.getCamsEdit(dataDB, dataRTSP)
			if len(camsForEdit) == 0 {
				a.log.Info("Data is identity, waiting...")
				continue
			}

			a.log.Info("Count of data is same, but the cameras are different")
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
				err := a.editCamerasToRTSP(ctx, camsForEdit)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}
			}
		}
	}
}
