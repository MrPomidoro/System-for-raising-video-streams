package main

import (
	"context"
	"fmt"
	"reflect"

	service "github.com/Kseniya-cha/System-for-raising-video-streams/cmd"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/logger"
)

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())

	// Чтение конфигурационного файла
	cfg, err := config.GetConfig()
	if err != nil || reflect.DeepEqual(cfg.Database, config.Database{}) {
		fmt.Println("cannot read config")
		return
	}

	log := logger.NewLogger(cfg)

	// Инициализация прототипа приложения
	app, err := service.NewApp(ctx, cfg)
	if err != nil {
		log.Error(err.Error())
		return
	}
	appI := service.App(app)

	// Запуск алгоритма в отдельной горутине
	go appI.Run(ctx)

	// Ожидание прерывающего сигнала
	appI.GracefulShutdown(cancel)
}
