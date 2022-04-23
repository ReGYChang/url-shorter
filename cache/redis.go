package cache

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"time"
	"url-shorter/helper"
)

var (
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func Add(c *gin.Context) {
	target := c.PostForm("target")
	expire := c.PostForm("expire")
	_expire, err := strconv.Atoi(expire)

	id := getCounter()

	hash := helper.DecToAny(id)

	err = rdb.Set(hash, target, time.Duration(_expire)*time.Second).Err()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "JSON_SUCCESS",
		"message":   "ok",
		"short_url": fmt.Sprintf("https://%s/%s", c.Request.Host, hash),
	})
}

func Visit(c *gin.Context) {
	hash := c.Param("hash")
	url, err := rdb.Get(hash).Result()
	if err != nil || len(url) < 1 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "JSON_ERROR",
			"message": "Not found",
		})
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func getCounter() int {
	rdb.Incr("counter")
	id, _ := rdb.Get("counter").Result()
	_id, _ := strconv.Atoi(id)
	return _id
}
