package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := ConnectToDB()
	if err != nil {
		panic(err)
	}
	err = InitDB(db)
	if err != nil {
		panic(err)
	}

	port := flag.Int("port", 8080, "Port to run on")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.GET("/get/*website", getWebsite)

	err = r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}

func getWebsite(ctx *gin.Context) {
	website := ctx.Param("website")
	ctx.JSON(http.StatusOK, gin.H{"website": website})
}
