package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/rtsp"
)

// Get запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) []refreshstream.RefreshStream {
	req, err := a.refreshStreamUseCase.Get(ctx)
	if err != nil {
		logger.LogError(a.Log, fmt.Sprintf("cannot get response from database: %v", err))
	}
	return req
}

/*
Получение спискоа камер с бд и с rtsp
На выходе: данные с бд, данные с rtsp, длины этих списков, статус код, ошибка
*/
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream, map[string]interface{}, int, int, string, error) {
	var lenResRTSP int

	// Отправка запросов к базе и к rtsp
	resDB := a.getReqFromDB(ctx)
	resRTSP := rtsp.GetRtsp(a.cfg)

	// resDB = []refreshstream.RefreshStream{} // проверка нулевого ответа от базы
	// Проверка, что ответ от базы данных не пустой
	if len(resDB) == 0 {
		return resDB, resRTSP, len(resDB), lenResRTSP, "400", errors.New("response from database is null")
	}

	// Определение числа потоков с rtsp
	for _, items := range resRTSP { // items - поле "items"
		// мапа: ключ - номер камеры, значения - остальные поля этой камеры
		camsMap := items.(map[string]interface{})
		lenResRTSP = len(camsMap) // количество камер
	}

	// Проверка, что ответ от rtsp данных не пустой
	if lenResRTSP == 0 {
		return resDB, resRTSP, len(resDB), lenResRTSP, "500", errors.New("response from rtsp-simple-server is null")
	}

	return resDB, resRTSP, len(resDB), lenResRTSP, "200", nil
}

/*
Функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая true, если количество камер в базе и в rtsp одинаковое
*/
func CheckIdentity(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) bool {

	// Счётчик для подсчёта совпадающих стримов камер
	var count int

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
					count++
					break
				}
			}
		}
	}

	// Если счётчик равен длине списка с базы данных, данные совпадают
	if count == len(dataDB) {
		return true
	} else {
		return false
	}
}

/*
Функция, принимающая на вход результат выполнения get запроса к базе и запроса к rtsp,
возвращающая списки отличающихся камер:
камеры, отсутствующие в rtsp, но имеющиеся в базе, добавить в rtsp,
лишние - удалить из rtsp
*/
func GetDifferenceElements(dataDB []refreshstream.RefreshStream, dataRTSP map[string]interface{}) ([]string, []string) {
	/*
		Слайсы с отличающимися элементами:
		лишние камеры из rtsp нужно удалить,
		отсутствующие в rtsp, но имеющиеся в базе - добавить
	*/
	var resSliceRemove []string
	var resSliceAdd []string

	// Перебор элементов списка структур
	for _, camDB := range dataDB {
		// Счетчик для проверки вхождения камеры бд в список с rtsp
		var isUniqueDB int

		// Перебор элементов мапы
		for _, camsRTSP := range dataRTSP {

			// Для возможности извлечения данных
			camsRTSPMap := camsRTSP.(map[string]interface{})

			// Счетчик для проверки вхождения камеры rtsp в список с бд
			var isUniqueRTSP int
			// Переменная для фиксации стрима камеры
			var cam string

			// camRTSP - стрим камеры
			for camRTSP := range camsRTSPMap {
				cam = camRTSP
				// Если stream из базы данных совпадает с rtsp, счётчики увеличиваются
				if camDB.Stream.String == camRTSP {
					isUniqueRTSP++
					isUniqueDB++
					break
				}
			}
			// Если значение счётчика ненулевое, камера добавляется в список на удаление
			if isUniqueRTSP == 0 {
				resSliceRemove = append(resSliceRemove, cam)
			}
		}
		// Если значение счётчика ненулевое, камера добавляется в список на добавление
		if isUniqueDB == 0 {
			resSliceAdd = append(resSliceAdd, camDB.Stream.String)
		}
	}

	return resSliceAdd, resSliceRemove
}

func EqualData() error {
	fmt.Println("run EqualData")
	return nil
}

func LessData() error {
	fmt.Println("run LessData")
	return nil
}

func MoreData() error {
	fmt.Println("run MoreData")
	return nil
}

// Вывод списка потоков с rtsp-simple-server
// (потом будет удалена или изменена, сейчас помогает разобраться)
func (a *app) getReqFromRtsp() {
	rtspResultMap := rtsp.GetRtsp(a.cfg)

	for key, items := range rtspResultMap { // items - значение поля "items"
		fmt.Printf("%T\n", rtspResultMap[key])
		// Для доступа к данным каждой камеры:
		camsMap := items.(map[string]interface{})

		for _, camFields := range camsMap { //
			// fmt.Println(camName) // camName - номер каждой камеры
			// Для доступа к данным полей камеры:
			camFieldsMap := camFields.(map[string]interface{})
			for fields, _ := range camFieldsMap { //valOfFields
				fmt.Println(fields) // поля confName, conf, source, sourceReady, tracks, readers
				// fmt.Printf("type: %T; value: %v\n\n", valOfFields, valOfFields) // значения этих полей
			}
		}
	}
}
