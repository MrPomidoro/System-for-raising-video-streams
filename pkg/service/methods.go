package service

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/rtsp"
)

// Get-запрос на получение списка камер из базы данных
func (a *app) getReqFromBD(ctx context.Context) {
	req, err := a.refreshStreamUseCase.Get(ctx)
	if err != nil {
		logger.LogError(a.Log, fmt.Sprintf("cannot get response from database: %v", err))
	} else {
		logger.LogDebug(a.Log, fmt.Sprintf("Response from database:\n%v", req))
	}
}

// Получение и парсинг списка потоков с rtsp-simple-server
func (a *app) getReqFromRtsp() {
	rtspResult := rtsp.GetRtsp(a.cfg)
	// Для доступа к полученным данным:
	rtspResultMap := rtspResult.(map[string]interface{})
	for _, items := range rtspResultMap { // items - поле "items"
		// Для доступа к данным каждой камеры:
		camsMap := items.(map[string]interface{})
		for camName, camFields := range camsMap {
			fmt.Println(camName) // camName - номер камеры
			// Для доступа к данным полей камеры:
			camFieldsMap := camFields.(map[string]interface{})
			for fields, valOfFields := range camFieldsMap {
				fmt.Println(fields)                                             // поля confName, conf, source, sourceReady, tracks, readers
				fmt.Printf("type: %T; value: %v\n\n", valOfFields, valOfFields) // значения этих полей
			}
			fmt.Printf("\n\n")
		}
	}
}
