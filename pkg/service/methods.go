package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/internal/refreshstream"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/rtsp"
)

// Get-запрос на получение списка камер из базы данных
func (a *app) getReqFromDB(ctx context.Context) []refreshstream.RefreshStream {
	req, err := a.refreshStreamUseCase.Get(ctx)
	if err != nil {
		logger.LogError(a.Log, fmt.Sprintf("cannot get response from database: %v", err))
	} else {
		logger.LogDebug(a.Log, fmt.Sprintf("Response from database:\n%T", req))
	}
	return req
}

// Вывод списка потоков с rtsp-simple-server
// (потом будет удалена, сейчас помогает разобраться)
func (a *app) getReqFromRtsp() {
	var i int
	rtspResultMap := rtsp.GetRtsp(a.cfg)

	for _, items := range rtspResultMap { // items - поле "items"
		// Для доступа к данным каждой камеры:
		camsMap := items.(map[string]interface{})
		fmt.Println(len(camsMap))

		for camName, _ := range camsMap { // camFields
			i++
			fmt.Println(camName) // camName - номер камеры
			// Для доступа к данным полей камеры:
			// camFieldsMap := camFields.(map[string]interface{})
			// for fields, _ := range camFieldsMap { //valOfFields
			// 	fmt.Println(fields) // поля confName, conf, source, sourceReady, tracks, readers
			// fmt.Printf("type: %T; value: %v\n\n", valOfFields, valOfFields) // значения этих полей
			// }
		}
	}
}

func (a *app) comparisonResultsRtspDb(ctx context.Context) {

}
