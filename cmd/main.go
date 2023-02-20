package main

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	ce "github.com/Kseniya-cha/System-for-raising-video-streams/pkg/customError"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/service"
)

// /*
func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())
	cerr := ce.NewError(ce.FatalLevel, "50.0.0", "error at main package level")

	// Чтение конфигурационного файла
	cfg, err := config.GetConfig()
	if err != nil {
		cerr.NextError(err)
		fmt.Println(cerr.SetError(err).Error())
		return
	}

	// log := logger.NewLogger(cfg)

	// Инициализация прототипа приложения
	app, err := service.NewApp(ctx, cfg)
	if err != nil {
		cerr.NextError(err)
		fmt.Println(cerr.Error())
		// fmt.Println(cerr.SetError().Error())
		return
	}

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Ожидание прерывающего сигнала
	// app.GracefulShutdown(app.SigChan)
	app.GracefulShutdown(cancel)
}

// */
