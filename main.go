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
	config, err := config.FromFile("./configuration.json")
	if err != nil {
		log.Fatal(err)
	}

	service, err := redis.New(config.Redis.Host, config.Redis.Port, config.Redis.Password)
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

	routersInit := handler.New(config.Options.Schema, config.Options.Prefix, service)
	readTimeout := config.Server.ReadTimeout
	writeTimeout := config.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%d", config.Server.Port)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()

	//router.Run(":"+configuration.Server.Port, router.Handler)

	//r := gin.Default()
	//
	//r.POST("/", storage.Save)
	//r.GET("/:hash", storage.Visit)
	//
	//r.Run(":9099")
}
