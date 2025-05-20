package main

import (
	"database/sql"
	"datadb"
	"errors"
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
	r.POST("/embed-text", handler(embedText))
	r.POST("/import-url", importURLHandler)
	r.POST("/semantic-search", handler(semanticSearch))
	r.POST("/agent-ramsay", handler(agentRamsaySearch))

	fmt.Println("Starting import server")
	err = r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func handler[T, U any](f func(T) (U, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := *new(T)
		err := ctx.BindJSON(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": errors.Join(errors.New("invalid body format"), err).Error(),
			})
			return
		}
		response, err := f(request)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.Join(errors.New("an error was encountered while processing request"), err).Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
