package main

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/service"
)

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())

	// Чтение конфигурационного файла
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	log := logger.NewLogger(cfg)
	// Инициализация прототипа приложения
	app, err := service.NewApp(ctx, cfg)
	if err != nil {
		log.Error(ce.NewError(ce.FatalLevel, "50.0.0", "error at main package level").Error())
	}

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Ожидание прерывающего сигнала
	// app.GracefulShutdown(app.SigChan)
	app.GracefulShutdown(ctx, cancel)
}
