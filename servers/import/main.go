package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("port", 8080, "Port to run on")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheck)

	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
