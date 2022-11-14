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

// Получение спискоа камер с бд и с rtsp
func (a *app) getDBAndApi(ctx context.Context) ([]refreshstream.RefreshStream, map[string]interface{}, int, int, string, error) {
	var lenResRTSP int

	// Отправка запросов к базе и к rtsp
	resDB := a.getReqFromDB(ctx)
	// a.getReqFromRtsp()
	resRTSP := rtsp.GetRtsp(a.cfg)

	resDB = []refreshstream.RefreshStream{} // проверка нулевого ответа от базы
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
		return a.getReqFromDB(ctx), resRTSP, len(resDB), lenResRTSP, "500", errors.New("response from rtsp-simple-server is null")
	}

	return a.getReqFromDB(ctx), resRTSP, len(resDB), lenResRTSP, "200", nil
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
