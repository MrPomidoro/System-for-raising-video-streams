package main

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/service"
)

type ctxLogger struct{}

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())

	// Чтение конфигурационного файла
	cfg, err := config.ConfigI.GetConfig(config.NewConfig())
	log := logger.NewLogger(cfg)
	if err != nil {
		log.Error(err.Error())
	}
	ctx = context.WithValue(ctx, ctxLogger{}, log)

	// Инициализация прототипа приложения
	app := service.NewApp(ctx, cfg)

	// Запуск алгоритма в отдельной горутине
	go service.App.Run(app, ctx)

	// Ожидание прерывающего сигнала
	// app.GracefulShutdown(app.SigChan)
	service.App.GracefulShutdown(app, ctx, cancel)
}
