package methods

import (
	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
)

/*
Функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая true, если количество камер в базе и в rtsp одинаковое
*/
func CheckIdentity(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) (bool, bool, []rtspsimpleserver.Conf) {

	var conf rtspsimpleserver.Conf
	var confArr []rtspsimpleserver.Conf

	// Счётчик для подсчёта совпадающих стримов камер
	var count int
	// Счётчик для подсчёта камер, у которых совпадают поля
	var identity int

	// Перебор элементов списка структур
	for _, camDB := range dataDB {

		// Перебор элементов мапы
		for _, camsRTSP := range dataRTSP {

			// Для возможности извлечения данных
			camsRTSPMap := camsRTSP.(map[string]interface{})

			// camStreamRTSP - стрим камеры
			for camStreamRTSP := range camsRTSPMap {

				// Если stream из базы данных не совпадает с rtsp, итерация пропускается
				if camDB.Stream.String != camStreamRTSP {
					continue
				}
				// Если совпадает - увеличивается счётчик количества совпадающих стримов
				count++

				camFields := camsRTSPMap[camStreamRTSP]            // все поля камеры (conf, confName, source etc)
				camFieldsMap := camFields.(map[string]interface{}) // для извлечения данных

				// camFieldName - имя поля ("sourceProtocol"), camField - значение ("tcp")
				for camFieldName, camField := range camFieldsMap {
					// рассматриваем только поле conf
					if camFieldName != "conf" {
						continue
					}

					camFieldMap := camField.(map[string]interface{}) // для извлечения данных

					conf.Stream = camStreamRTSP
					// fmt.Printf("camStreamRTSP=%s, camDB.Stream.String=%s\n", camStreamRTSP, camDB.Stream.String)
					// fmt.Printf("camFieldMap[sourceProtocol].(string) = %s, camDB.Protocol.String=%s\n", camFieldMap["sourceProtocol"].(string), camDB.Protocol.String)

					// Если значение поля в rtsp отличается от значения в бд, данные из бд вносятся в структуру
					if camFieldMap["sourceProtocol"].(string) != camDB.Protocol.String {
						conf.SourceProtocol = camDB.Protocol.String
						continue
					}
					identity++
					// fmt.Println("identity int:", identity)
					// парсинг поля runOnReady
					// var runOnReadyRes string
					// runOnReady := camFieldMap["runOnReady"].(string)
					// runOnReadyArr := strings.Split(runOnReady, " --")
					// for i, elemRunOnReadyArr := range runOnReadyArr {
					// 	elemRunOnReadyArrArr := strings.Split(elemRunOnReadyArr, " ")
					// 	if elemRunOnReadyArrArr[0] == "port" {
					// 		if elemRunOnReadyArrArr[1] != camDB.Portsrv {
					// 			runOnReadyRes
					// 		}
					// 	}
					// }

				}
				break

			}
			confArr = append(confArr, conf)
		}
	}

	// Если счётчик равен длине списка с базы данных, данные совпадают
	if count != len(dataDB) {
		return false, false, confArr
	} else if count == len(dataDB) && identity == len(dataDB) {
		return true, true, confArr
	} else if count == len(dataDB) && identity != len(dataDB) {
		return true, false, confArr
	}
	return false, false, confArr
}

/*
Функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая списки отличающихся камер:
камеры, отсутствующие в rtsp, но имеющиеся в базе, нужно добавить в rtsp,
камеры, имеющиеся в rtsp, но отсутствующие в базе, - удалить из rtsp
*/
func GetDifferenceElements(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) ([]string, []string) {

	// Слайсы с отличающимися элементами:
	// лишние камеры из rtsp нужно удалить,
	// отсутствующие в rtsp, но имеющиеся в базе - добавить
	var resSliceRemove []string
	var resSliceAdd []string

	// Счётчики
	var doubleRemove int
	var doubleAppend int

	// Формирование списка на удаление камер из rtsp

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

	// Формирование списка на добавление камер в rtsp

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

		// Если значение счётчика ненулевое, камера добавляется в список на добавление
		if doubleAppend == 0 {
			resSliceAdd = append(resSliceAdd, camDB.Stream.String)
		}
		doubleAppend = 0
	}

	return resSliceAdd, resSliceRemove
}
