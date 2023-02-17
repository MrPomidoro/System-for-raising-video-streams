package main

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/service"
)

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())

	// Чтение конфигурационного файла
	cfg, _ := config.GetConfig()

	// Инициализация прототипа приложения
	app := service.NewApp(ctx, cfg)

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Ожидание прерывающего сигнала
	// app.GracefulShutdown(app.SigChan)
	app.GracefulShutdown(ctx, cancel)
}
