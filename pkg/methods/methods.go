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
func CheckIdentityAndCountOfData(dataDB []refreshstream.RefreshStream, dataRTSP []rtspsimpleserver.SConf, cfg *config.Config) (bool, bool, []rtspsimpleserver.SConf) {

	var sconfArr []rtspsimpleserver.SConf

	// Счётчик для подсчёта совпадающих стримов камер
	var count int
	// Счётчик для подсчёта камер, у которых совпадают поля
	var identity int

	// Перебор элементов списка структур
	for _, camDB := range dataDB {
		var sconf rtspsimpleserver.SConf

		// Перебор элементов мапы, _ - "items"
		for _, camRTSP := range dataRTSP {

			// Если stream из базы данных не совпадает с rtsp, итерация пропускается
			if camDB.Stream.String != camRTSP.Stream {
				continue
			}
			// Если совпадает - увеличивается счётчик количества совпадающих стримов
			count++

			// ------------------------------------------- //
			//   Проверка одинаковости данных для камеры   //
			// ------------------------------------------- //

			sconf.Stream = camRTSP.Stream

			if camRTSP.Conf.SourceProtocol == camDB.Protocol.String {
				identity++
			} else {
				sconf.Conf.SourceProtocol = camDB.Protocol.String
			}

			// парсинг поля runOnReady
			var runOnReady string
			// sconf.Conf.RunOnReady = runOnReady

			if cfg.Run == "" {
				if camRTSP.Conf.RunOnReady == "" {
					// if camFieldMap["runOnReady"].(string) == "" {
					identity++
				}
			} else {
				runOnReady = fmt.Sprintf(cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)
				if camRTSP.Conf.RunOnReady == runOnReady {
					identity++
				} else {
					sconf.Conf.RunOnReady = runOnReady
				}
			}

			// парсинг поля source
			var source = fmt.Sprintf("rtsp://%s@%s/%s", camDB.Auth.String, camDB.Ip.String, camDB.Stream.String)

			if camRTSP.Conf.Source == source {
				// if camFieldMap["source"].(string) == source {
				identity++
				// continue
			} else {
				sconf.Conf.Source = source
			}
			break
		}
		sconfArr = append(sconfArr, sconf)
	}
	// }

	lenDB := len(dataDB)
	// Если счётчик равен длине списка с базы данных, данные совпадают
	countEqual, identityEqual := compareDBandRTSP(count, identity, lenDB)
	return countEqual, identityEqual, sconfArr
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

// GetLensData проверяет, что полученные ответы от rtsp и базы не пустые, и возвращает их длины
func GetLensData(resDB []refreshstream.RefreshStream, resRTSP map[string]interface{}) (int, int) {
	var lenResRTSP int

	// Определение числа потоков с rtsp
	for _, items := range resRTSP { // items - поле "items"
		// мапа: ключ - номер камеры, значения - остальные поля этой камеры
		camsMap := items.(map[string]interface{})
		lenResRTSP = len(camsMap) // количество камер
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
func GetCamsForRemove(dataDB []refreshstream.RefreshStream, dataRTSP []rtspsimpleserver.SConf) []string {

	// Слайс с камерами, отсутствующим в rtsp, но имеющимися в базе
	var resSliceRemove []string
	// Счётчик
	var doubleRemove int

	// Перебор элементов списка структур
	for _, camDB := range dataDB {
		// Перебор элементов мапы
		for _, camRTSP := range dataRTSP {
			if camDB.Stream.String == camRTSP.Stream {
				doubleRemove++
				break
			}

			if camDB.Stream.String == camRTSP.Stream {
				doubleRemove++
				break
			}

			// Если значение счётчика ненулевое, камера добавляется в список на удаление
			if doubleRemove == 0 {
				resSliceRemove = append(resSliceRemove, camRTSP.Stream)
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
func GetCamsForAdd(dataDB []refreshstream.RefreshStream, dataRTSP []rtspsimpleserver.SConf) []string {

	// Слайс с камерами, отсутствующими в rtsp, но имеющимися в базе
	var resSliceAdd []string
	// Счётчик
	var doubleAppend int

	// Перебор элементов списка структур
	for _, camDB := range dataDB {

		for _, camRTSP := range dataRTSP {
			if camRTSP.Stream == camDB.Stream.String {
				doubleAppend++
				break
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
