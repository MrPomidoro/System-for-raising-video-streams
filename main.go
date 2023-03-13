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
		fmt.Println("ERROR: cannot read config: file is empty")
		return
	}

	log := logger.NewLogger(cfg)

	// Инициализация прототипа приложения
	app, err := service.NewApp(ctx, cfg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Ожидание прерывающего сигнала
	app.GracefulShutdown(cancel)
}
