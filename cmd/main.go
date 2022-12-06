package main

import (
	"context"

	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/config"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/database"
	"github.com/Kseniya-cha/System-for-raising-video-streams/pkg/service"
)

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())

	// Чтение конфигурационного файла
	cfg := config.GetConfig()

	// Инициализация прототипа приложения
	app := service.NewApp(cfg)

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Проверка коннекта к базе данных
	// и переподключение при необходимости
	go database.DBPing(cfg, app.Db, ctx)

	// Ожидание прерывающего сигнала
	// app.GracefulShutdown(app.SigChan)
	app.GracefulShutdown(app.SigChan, cancel)
}
