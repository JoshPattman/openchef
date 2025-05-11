package main

import (
	"database/sql"
	"datadb"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"utils"

	"github.com/gin-gonic/gin"
)

var logger *slog.Logger
var db *sql.DB

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	var err error
	db, err = ConnectToDB()
	if err != nil {
		panic(err)
	}
	err = datadb.InitDB(db)
	if err != nil {
		panic(err)
	}

	port := utils.MustReadEnvInt("DATA_PORT")

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/healthcheck", healthCheckHandler)
	r.POST("/embed-text", embedTextHandler)
	r.POST("/import-url", importURLHandler)
	r.POST("/semantic-search", semanticSearchHandler)

	fmt.Println("Starting import server")
	err = r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
