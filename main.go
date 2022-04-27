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
	cfg, err := config.FromFile("./configuration.json")
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

	//gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := handler.New(cfg.Options.Schema, cfg.Options.Prefix, service)
	readTimeout := cfg.Server.ReadTimeout
	writeTimeout := cfg.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%d", cfg.Server.Port)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server startup err", err)
		}
	}()

	//router.Run(":"+configuration.Server.Port, router.Handler)

	//r := gin.Default()
	//
	//r.POST("/", storage.Save)
	//r.GET("/:hash", storage.Visit)
	//
	//r.Run(":9099")
}
