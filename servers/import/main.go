package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"utils"

	"github.com/gin-gonic/gin"
)

var logger *slog.Logger

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	port := utils.MustReadEnvInt("IMPORT_PORT")

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.POST("/basic-info", basicInfoHandler)
	r.POST("/advanced-info", advancedInfoHandler)
	r.POST("/embed-text", embedTextHandler)

	fmt.Println("Starting import server")
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
