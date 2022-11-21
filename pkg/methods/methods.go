package methods

import (
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	rtspsimpleserver "github.com/Kseniya-cha/System-for-raising-video-streams/internal/rtsp-simple-server"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
)

/*
Функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая true, если количество камер в базе и в rtsp одинаковое
*/
func CheckIdentity(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}, cfg *config.Config) (bool, bool, []rtspsimpleserver.Conf) {

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

				// --------------------------------------- //
				// Проверка одинаковости данных для камеры //
				// --------------------------------------- //

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
						conf.SourceProtocol = camFieldMap["sourceProtocol"].(string)
					}

					// парсинг поля runOnReady
					runOnReady := fmt.Sprintf(cfg.Run, camDB.Portsrv, camDB.Sp.String, camDB.CamId.String)

					if camFieldMap["runOnReady"].(string) == runOnReady {
						identity++
						continue
					}
					conf.RunOnReady = runOnReady

				}
				break

			}
			confArr = append(confArr, conf)
		}
	}

	// Если счётчик равен длине списка с базы данных, данные совпадают
	if count != len(dataDB) {
		return false, false, confArr
	} else if count == len(dataDB) && identity == 2*len(dataDB) {
		return true, true, confArr
	} else if count == len(dataDB) && identity != 2*len(dataDB) {
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
