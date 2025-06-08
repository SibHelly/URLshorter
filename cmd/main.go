package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/SibHelly/url_shortnener/internals/app"
	"github.com/SibHelly/url_shortnener/internals/cfg"
)

func main() {
	// загрузка конфигов
	config := cfg.LoadAndStoreConfig()

	ctx, cancel := context.WithCancel(context.Background())

	// получение сигналов системы
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//создание сервера
	server := app.NewServer(config, ctx)

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		server.Shutdown()
		cancel()

	}()

	// запуск сервера
	server.Serve()

}
