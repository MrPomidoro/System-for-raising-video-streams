package methods

import (
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
)

/*
CheckIdentityAndCountOfData - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp, возвращающая
  - (true, true), если количество камер в базе и в rtsp одинаковое;
  - (true, false), если количество камер в базе и в rtsp одинаковое, но сами данные отличаются;
  - (false, false), если количество камер в базе и в rtsp отличается.

Также возвращается список, содержащий поля sourceProtocol и runOnReady, если в rtsp они отличаются от данных из бд,
и стрим камеры (всегда)
*/
func CheckIdentityAndCountOfData(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}, cfg *config.Config) (bool, bool, []rtspsimpleserver.Conf) {

	var confArr []rtspsimpleserver.Conf

	// Счётчик для подсчёта совпадающих стримов камер
	var count int
	// Счётчик для подсчёта камер, у которых совпадают поля
	var identity int

	// Перебор элементов списка структур
	for _, camDB := range dataDB {

		// Перебор элементов мапы, _ - "items"
		for _, camsRTSP := range dataRTSP {

			// Для возможности извлечения данных
			camsRTSPMap := camsRTSP.(map[string]interface{})

			var conf rtspsimpleserver.Conf

			// camStreamRTSP - стрим камеры, camFields - все поля камеры (conf, confName, source etc)
			for camStreamRTSP, camFields := range camsRTSPMap {

				// Если stream из базы данных не совпадает с rtsp, итерация пропускается
				if camDB.Stream.String != camStreamRTSP {
					continue
				}
				// Если совпадает - увеличивается счётчик количества совпадающих стримов
				count++

				// ------------------------------------------- //
				//   Проверка одинаковости данных для камеры   //
				// ------------------------------------------- //

				camFieldsMap := camFields.(map[string]interface{}) // для извлечения данных

				// camFieldName - имя поля ("conf"), camField - значение (поля) этого поля ("sourceProtocol")
				for camFieldName, camField := range camFieldsMap {
					// рассматриваем только поле conf
					if camFieldName != "conf" {
						continue
					}

					camFieldMap := camField.(map[string]interface{}) // для извлечения данных

					conf.Stream = camStreamRTSP

					// Если значение поля в rtsp отличается от значения в бд, данные из бд вносятся в структуру
					if camFieldMap["sourceProtocol"].(string) == camDB.Protocol.String {
						identity++
					} else {
						conf.SourceProtocol = camDB.Protocol.String
					}

					// парсинг поля runOnReady
					var runOnReady string
					conf.RunOnReady = runOnReady

					if cfg.Run == "" {
						if camFieldMap["runOnReady"].(string) == "" {
							identity++
						}
					} else {
						runOnReady = fmt.Sprintf(cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)
						if camFieldMap["runOnReady"].(string) == runOnReady {
							identity++
						} else {
							conf.RunOnReady = runOnReady
						}
					}

					// парсинг поля source
					var source = fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream.String)

					if camFieldMap["source"].(string) == source {
						identity++
						continue
					} else {
						conf.Source = source
					}
					// fmt.Printf("methods: conf - %#v\n\n", conf)

				}
				break

			}
			confArr = append(confArr, conf)
		}
	}

	lenDB := len(dataDB)
	// Если счётчик равен длине списка с базы данных, данные совпадают
	countEqual, identityEqual := compareDBandRTSP(count, identity, lenDB)
	return countEqual, identityEqual, confArr
}

// compareDBandRTSP сравнивает счётчик count и длине списка с базы данных,
// а также счетчик identity с утроенной длиной списка с базы данных
func compareDBandRTSP(count, identity, lenDB int) (bool, bool) {
	if count != lenDB {
		return false, false
	} else if count == lenDB && identity == 3*lenDB {
		return true, true
	} else if count == lenDB && identity != 3*lenDB {
		return true, false
	}
	return false, false
}

//

// checkEmptyData проверяет, что полученные ответы от rtsp и базы не пустые, и возвращает их длины
func CheckEmptyData(resDB []refreshstream.RefreshStream, resRTSP map[string]interface{}) (int, int) {
	var lenResRTSP int

	// Проверка, что ответ от базы данных не пустой
	if len(resDB) == 0 {
		return 0, 0
	}

	// Определение числа потоков с rtsp
	for _, items := range resRTSP { // items - поле "items"
		// мапа: ключ - номер камеры, значения - остальные поля этой камеры
		camsMap := items.(map[string]interface{})
		lenResRTSP = len(camsMap) // количество камер
	}

	// Проверка, что ответ от rtsp данных не пустой
	if lenResRTSP == 0 {
		return 0, 0
	}

	return len(resDB), lenResRTSP
}

//
//
//

/*
GetCamsForRemove - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая список камер, имеющихся в rtsp, но отсутствующих в базе
*/
func GetCamsForRemove(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) []string {

	// Слайс с камерами, отсутствующим в rtsp, но имеющимися в базе
	var resSliceRemove []string
	// Счётчик
	var doubleRemove int

	// Перебор элементов мапы
	for _, camsRTSP := range dataRTSP {
		// Для возможности извлечения данных
		camsRTSPMap := camsRTSP.(map[string]interface{})
		// camRTSP - стрим камеры
		for camRTSP := range camsRTSPMap {
			// Перебор элементов списка структур
			for _, camDB := range dataDB {

				if camDB.Stream.String == camRTSP {
					doubleRemove++
					break
				}
			}

			// Если значение счётчика ненулевое, камера добавляется в список на удаление
			if doubleRemove == 0 {
				resSliceRemove = append(resSliceRemove, camRTSP)
			}
			doubleRemove = 0
		}
	}

	return resSliceRemove
}

/*
GetCamsForAdd - функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая список камер, отсутствующих в rtsp, но имеющихся в базе
*/
func GetCamsForAdd(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) []string {

	// Слайс с камерами, отсутствующим в rtsp, но имеющимися в базе
	var resSliceAdd []string
	// Счётчик
	var doubleAppend int

	// Перебор элементов списка структур
	for _, camDB := range dataDB {

		// Перебор элементов мапы
		for _, camsRTSP := range dataRTSP {
			// Для возможности извлечения данных
			camsRTSPMap := camsRTSP.(map[string]interface{})
			// camRTSP - стрим камеры
			for camRTSP := range camsRTSPMap {
				// Если stream из базы данных совпадает с rtsp, счётчики увеличиваются
				if camDB.Stream.String == camRTSP {
					doubleAppend++
					break
				}
			}
		}

		// Если значение счётчика ненулевое, камера попадает в список на добавление
		if doubleAppend == 0 {
			resSliceAdd = append(resSliceAdd, camDB.Stream.String)
		}
		doubleAppend = 0
	}

	return resSliceAdd
}
