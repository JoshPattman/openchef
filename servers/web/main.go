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

	port := flag.Int("port", 8080, "Port to run on")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.GET("/get/*website", getWebsite)

	err = r.Run(fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}

func getWebsite(ctx *gin.Context) {
	website := ctx.Param("website")

	request := utils.ImportFromURLRequest{
		URL: website,
	}

	importPort := flag.Lookup("import_port").Value.String()
	url := fmt.Sprintf("http://localhost:%s/basic_info", importPort)

	resp, err := http.Post(url, "application/json", utils.ToJSONReader(request))
	if err != nil {
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
