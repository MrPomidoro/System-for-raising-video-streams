package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/methods"
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
	// и переподключение при необходимост
	if a.db == nil {
		return
	}

	go a.db.DBPing(ctx, a.cfg)

loop:
	for {
		fmt.Println("")

		select {

		// Если контекст закрыт, loop завершается
		case <-ctx.Done():
			break loop

		// Выполняется периодически через установленный в конфигурационном файле промежуток времени
		case <-tick.C:
			fmt.Println("aaaaaaaaaaaa")
			// Получение данных от базы данных и от rtsp
			dataDB, dataRTSP, err := a.getDBAndApi(ctx)
			if err != nil {
				a.log.Error(err.Error())
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
			return
			if lenResDB == lenResRTSP {
				a.log.Info(fmt.Sprintf("The count of data in the database = %d is equal to the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))

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
					a.log.Error(err.Error())
					continue
				}

				//
				/*
					Если данных в базе больше, чем в rtsp:
					- получение списка отличий;
					- API на добавление в ртсп;
					- запись в status_stream
				*/
			} else if lenResDB > lenResRTSP {

				a.log.Info(fmt.Sprintf("The count of data in the database = %d is greater than the count of data in rtsp-simple-server = %d", lenResDB, lenResRTSP))
				err = a.addAndRemoveData(ctx, dataRTSP, dataDB)
				if err != nil {
					a.log.Error(err.Error())
					continue
				}

				//
				/*
					Если данных в базе меньше, чем в rtsp:
					- получение списка отличий;
					- API на добавление в ртсп;
					- запись в status_stream
				*/
			} else if lenResDB < lenResRTSP {
				a.log.Info(fmt.Sprintf("The count of data in the database = %d is less than the count of data in rtsp-simple-server = %d; waiting...", lenResDB, lenResRTSP))

				// Ожидание 5 секунд и повторный запрос данных с базы и с rtsp
				time.Sleep(time.Second * 5)
				dataDB, dataRTSP, err := a.getDBAndApi(ctx)
				if err != nil {
					a.log.Error(err.Error())
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
						a.log.Error(err.Error())
						continue
					}

				} else {

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
