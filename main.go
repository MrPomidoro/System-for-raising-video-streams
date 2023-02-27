package main

import (
	"context"
	"fmt"

	service "github.com/Kseniya-cha/System-for-raising-video-streams/cmd"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	cerror "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())
	cerr := cerror.NewError(cerror.FatalLevel, "", "error at main package level")

	// Чтение конфигурационного файла
	cfg, err := config.GetConfig()
	if err != nil {
		cerr.NextError(err)
		fmt.Println(cerr.SetError(err).Error())
		return
	}

	log := logger.NewLogger(cfg)

	// Инициализация прототипа приложения
	app, err := service.NewApp(ctx, cfg)
	if err != nil {
		cerr.NextError(err)
		log.Error(cerr.Error())
		return
	}

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Ожидание прерывающего сигнала
	// app.GracefulShutdown(app.SigChan)
	app.GracefulShutdown(cancel)
}
