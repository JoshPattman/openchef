package main

import (
	"context"
	"encoding/json"
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

	fmt.Println("Hitting: ", url)

	resp, err := http.Post(url, "application/json", utils.ToJSONReader(request))
	if err != nil {
		fmt.Printf("Error: %s", err)
		ctx.String(http.StatusInternalServerError, "Error calling basic info endpoint")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.String(http.StatusInternalServerError, "Basic info endpoint returned an error")
		return
	}

	var recipe utils.Recipe
	err = json.NewDecoder(resp.Body).Decode(&recipe)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error decoding response")
		return
	}

	err = RecipePage(recipe.Name).Render(context.Background(), ctx.Writer)
	if err != nil {
		panic(err)
	}
}
