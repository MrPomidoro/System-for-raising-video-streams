package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
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

			// lenResDB, lenResRTSP := methods.GetLensData(dataDB, dataRTSP)

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

				camsForEdit := a.getCamsEdit(a.cfg, dataDB, dataRTSP)
				// Проверка одинаковости данных по стримам
				if len(camsForEdit) == 0 {
					a.log.Info("Data is identity, waiting...")
					continue
				}

				a.log.Info("Count of data is same, but the values are different")
				err := a.editCamerasToRTSP(ctx, camsForEdit)
				if err != nil {
					a.log.Error(err.Error())
				}

				return // edit

				// Если число данных совпадает и данные одинаковые ИЛИ если число данных совпадает, но данные отличаются,
				// метод equalOrIdentityData возвращает true
				// eqId := a.equalOrIdentityData(ctx, isEqualCount, identity, confArr)
				// if eqId {
				// 	continue
				// }

				return // edit

				// Если число данных отличается, выполняется differentCount
				// err := a.differentCount(ctx, dataDB, dataRTSP)
				// if err != nil {
				// 	a.log.Error(err.Error())
				// 	continue
				// }

				//
			/*
				Если данных в базе больше, чем в rtsp:
				- получение списка отличий;
				- API на добавление в ртсп;
				- запись в status_stream
			*/
			case len(dataDB) > len(dataRTSP):

				a.log.Info(fmt.Sprintf("The count of data in the database = %d is greater than the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))
				time.Sleep(8 * time.Second)
				return // edit

				// err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
				// if err != nil {
				// 	a.log.Error(err.Error())
				// 	continue
				// }

				//
			/*
				Если данных в базе меньше, чем в rtsp:
				- получение списка отличий;
				- API на добавление в ртсп;
				- запись в status_stream
			*/
			case len(dataDB) < len(dataRTSP):
				a.log.Info(fmt.Sprintf("The count of data in the database = %d is less than the count of data in rtsp-simple-server = %d; waiting...", lenResDB, lenResRTSP))

				// Ожидание 5 секунд и повторный запрос данных с базы и с rtsp
				time.Sleep(time.Second * 5)
				// dataDB, dataRTSP, err := a.getDBAndApi(ctx, &mu)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}
				dataDB, dataRTSP, err := a.getDBAndApi(ctx, &mu)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}
				// lenResDBLESS, lenResRTSPLESS := methods.GetLensData(dataDB, dataRTSP)

				// /*
				// Сравнение числа записей в базе данных и записей в rtsp после нового запроса
				if len(dataDB) == len(dataRTSP) {

					// Проверка одинаковости данных по стримам
					isEqualCount, identity, confArr := methods.CheckIdentityAndCountOfData(dataDB, dataRTSP, a.cfg)

					return // edit

					// Если число данных совпадает и данные одинаковые ИЛИ если число данных совпадает, но данные отличаются,
					// метод equalOrIdentityData возвращает true
					eqId := a.equalOrIdentityData(ctx, isEqualCount, identity, confArr)
					if eqId {
						continue
					}

					return // edit

					// Если число данных отличается, выполняется differentCount
					// err := a.differentCount(ctx, dataDB, dataRTSP)
					// if err != nil {
					// 	a.log.Error(err.Error())
					// 	continue
					// }

					// } else {

					// 	return // edit
					// 	err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
					// 	if err != nil {
					// 		a.log.Error(err.Error())
					// 		continue
					// 	}

				}
				// */
			}
		}
	}
}
