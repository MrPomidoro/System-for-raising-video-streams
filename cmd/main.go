package main

import (
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/service"
)

func main() {
	// Чтение конфигурационного файла
	cfg := config.GetConfig()

	// Инициализация прототипа приложения
	app := service.NewApp(cfg)
	logger.LogDebug(app.Log, "app created")

	// Запуск алгоритма в отдельной горутине
	go func() {
		if err := app.Run(); err != nil {
			logger.LogFatal(app.Log, err)
		}
	}()

	// Проверка коннекта к базе данных
	// и переподключение при необходимости
	go database.DBPing(cfg, app.Db)

	// Ожидание прерывающего сигнала
	app.GracefulShutdown(app.SigChan)
}
