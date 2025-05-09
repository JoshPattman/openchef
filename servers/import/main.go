package main

import (
	"fmt"
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
)

func main() {
	port := utils.MustReadEnvInt("IMPORT_PORT")

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.POST("/basic-info", basicInfoHandler)
	r.POST("/advanced-info", advancedInfoHandler)

	fmt.Println("Starting import server")
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
