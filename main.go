package main

import (
	"fmt"
	"log"
	"net/http"
	"url-shorter/config"
	"url-shorter/handler"
	"url-shorter/storage"
	"url-shorter/storage/redis"
)

func main() {
	cfg, err := config.FromFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	service, err := redis.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer func(service storage.Service) {
		err := service.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(service)

	handlerInit := handler.New(cfg.Options.Schema, cfg.Options.Prefix, service)
	readTimeout := cfg.Server.ReadTimeout
	writeTimeout := cfg.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%s", cfg.Server.Port)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        handlerInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("http server startup err", err)
	}
}
