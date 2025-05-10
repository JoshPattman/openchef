package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"utils"

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

	port := utils.MustReadEnvInt("WEB_PORT")
	//importPort := utils.MustReadEnvInt("IMPORT_PORT")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.GET("/get/*website", getWebsite)

	fmt.Println("Starting web server")
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
	if len(website) > 0 && website[0] == '/' {
		website = website[1:]
	}

	fmt.Println("Fetching website: ", website)

	request := utils.ImportFromURLRequest{
		URL: website,
	}

	importPort := utils.MustReadEnvInt("IMPORT_PORT")
	url := fmt.Sprintf("http://import_service:%d/basic-info", importPort)
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
