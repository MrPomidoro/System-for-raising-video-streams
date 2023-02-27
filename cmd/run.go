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

	ctx, _ = context.WithCancel(ctx)

	// Канал для периодического выполнения алгоритма
	tick := time.NewTicker(a.cfg.RefreshTime)
	defer tick.Stop()

	// Проверка коннекта к базе данных
	if a.db == nil {
		return
	}

	// Переподключение при необходимости
	go a.db.DBPing(ctx, a.cfg)

	var mu sync.Mutex

loop:
	for {
		select {

		case <-ctx.Done():
			break loop

		// Выполняется периодически через установленный в конфигурационном файле промежуток времени
		case <-tick.C:

			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, err := a.getDBAndApi(ctx, &mu)
			if err != nil {
				a.log.Error(err.Error())
				continue
			}

			// ---------------------------------------------------------- //
			//   Сравнение числа записей в базе данных и записей в rtsp   //
			// ---------------------------------------------------------- //

			switch {
			/*
				Если данных в базе столько же, сколько в rtsp:
				проверка, одинаковые ли записи:
				- если одинаковые, завершение и ожидание следующего запуска программы;
				- если различаются:
					- получение списка отличий,
					- отправка API,
					- запись в status_stream.
			*/
			case len(dataDB) == len(dataRTSP):
				a.log.Info(fmt.Sprintf("The count of data in the database = %d is equal to the count of data in rtsp-simple-server = %d", len(dataDB), len(dataRTSP)))

				// Получение отличающихся камер
				camsForEdit := a.getCamsEdit(a.cfg, dataDB, dataRTSP)
				if len(camsForEdit) == 0 {
					a.log.Info("Data is identity, waiting...")
					continue
				}

				// Если имеются отличия, отправляется запрос к ртсп на изменение
				a.log.Info("Count of data is same, but the values are different")
				err := a.editCamerasToRTSP(ctx, camsForEdit)
				if err != nil {
					a.log.Error(err.Error())
				}

			/*
				Если данных в базе больше, чем в rtsp:
				- получение списка отличий;
				- API на добавление в ртсп;
				- запись в status_stream
			*/
			case len(dataDB) > len(dataRTSP):

				a.log.Info(fmt.Sprintf("The count of data in the database = %d is greater than the count of data in rtsp-simple-server = %d", len(dataDB), len(dataRTSP)))
				// time.Sleep(5 * time.Second)

				err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}

			/*
				Если данных в базе меньше, чем в rtsp:
				- получение списка отличий;
				- API на добавление в ртсп;
				- запись в status_stream
			*/
			case len(dataDB) < len(dataRTSP):
				a.log.Info(fmt.Sprintf("The count of data in the database = %d is less than the count of data in rtsp-simple-server = %d; waiting...", len(dataDB), len(dataRTSP)))

				// Ожидание 5 секунд и повторный запрос данных с базы и с rtsp
				time.Sleep(time.Second * 5)

				dataDB, dataRTSP, err := a.getDBAndApi(ctx, &mu)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}

				// Сравнение числа записей в базе данных и записей в rtsp после нового запроса
				switch {
				case len(dataDB) == len(dataRTSP):

					a.log.Info(fmt.Sprintf("The count of data in the database = %d is equal to the count of data in rtsp-simple-server = %d", len(dataDB), len(dataRTSP)))

					// Получение отличающихся камер
					camsForEdit := a.getCamsEdit(a.cfg, dataDB, dataRTSP)
					if len(camsForEdit) == 0 {
						a.log.Info("Data is identity, waiting...")
						continue
					}

					// Если имеются отличия, отправляется запрос к ртсп на изменение
					a.log.Info("Count of data is same, but the values are different")
					err := a.editCamerasToRTSP(ctx, camsForEdit)
					if err != nil {
						a.log.Error(err.Error())
					}

				default:
					// Если число данных отличается, выполняется differentCount
					err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
					if err != nil {
						a.log.Error(err.Error())
						continue
					}

				}
			}
		}
	}
}
