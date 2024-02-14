package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the application
//
//	@title						Product API [capstone_backend]
//	@version					1.0
//	@description				This is swagger documentation for Product API
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
func main() {
	slogHandler := slog.Handler(slog.NewTextHandler(os.Stdout, nil))
	if !conf.Project.Debug {
		slogHandler = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(slogHandler).With("svc", conf.Project.ServiceName)

	logger.Info("Starting the application")
	
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
