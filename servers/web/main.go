package main

import (
	"context"
	"flag"
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

	port := utils.MustReadEnvInt("WEB_PORT")
	//importPort := utils.MustReadEnvInt("IMPORT_PORT")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.GET("/get/*website", getWebsite)

	fmt.Println("Starting web server")
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	logger.Info("Responding to healthcheck")
	ctx.JSON(http.StatusOK, "OK")
}

func getWebsite(ctx *gin.Context) {
	website := ctx.Param("website")
	if len(website) > 0 && website[0] == '/' {
		website = website[1:]
	}

	logger.Info("Fetching website: ", "website", website)

	request := utils.ImportFromURLRequest{
		URL: website,
	}

	importPort := utils.MustReadEnvInt("DATA_PORT")
	url := fmt.Sprintf("http://data_service:%d/basic-info", importPort)
	recipe := utils.Recipe{}
	err := PostAndReciveJson(url, request, &recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = RecipePage(recipe.Name).Render(context.Background(), ctx.Writer)
	if err != nil {
		panic(err)
	}
}
