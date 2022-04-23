package main

import (
	"github.com/gin-gonic/gin"
	"url-shorter/cache"
)

func main() {
	r := gin.Default()

	r.POST("/", cache.Add)
	r.GET("/:hash", cache.Visit)

	r.Run(":9099")
}
