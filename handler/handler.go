package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"time"
	"url-shorter/storage"
)

func New(schema string, host string, storage storage.Service) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	h := handler{schema, host, storage}
	r.POST("/encode", responseHandler(h.encode))
	r.GET("/:shortLink", h.redirect)
	r.GET("/info/:shortLink", responseHandler(h.decode))

	return r
}

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"shortUrl"`
}

type handler struct {
	schema  string
	host    string
	storage storage.Service
}

type input struct {
	URL     string `json:"url"`
	Expires string `json:"expires"`
}

func responseHandler(h func(ctx *gin.Context) (interface{}, int, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, status, err := h(ctx)
		if err != nil {
			data = err.Error()
		}
		ctx.JSON(status, response{Data: data, Success: err == nil})
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}
	}
}

func (h handler) encode(ctx *gin.Context) (interface{}, int, error) {
	input := input{}
	err := ctx.BindJSON(&input)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err)
	}

	uri, err := url.ParseRequestURI(input.URL)

	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid url")
	}

	layoutISO := "2006-01-02 15:04:05"
	expires, err := time.Parse(layoutISO, input.Expires)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid expiration date")
	}

	c, err := h.storage.Save(uri.String(), expires)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Could not store in database: %v", err)
	}

	u := url.URL{
		Scheme: h.schema,
		Host:   h.host,
		Path:   c,
	}

	fmt.Printf("Generated link: %v \n", u.String())

	return u.String(), http.StatusCreated, nil
}

func (h handler) decode(ctx *gin.Context) (interface{}, int, error) {
	code := ctx.Param("shortLink")

	model, err := h.storage.LoadInfo(code)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}

	return model, http.StatusOK, nil
}

func (h handler) redirect(ctx *gin.Context) {
	code := ctx.Param("shortLink")

	uri, err := h.storage.Load(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "JSON_ERROR",
			"message": "Data Not found",
		})
	}

	ctx.Redirect(http.StatusMovedPermanently, uri)
}
