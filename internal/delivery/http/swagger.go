package http

import (
	"fmt"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/ibookerke/capstone_backend/docs"
	"github.com/ibookerke/capstone_backend/internal/config"
)

type SwaggerHandler struct {
	conf config.Swagger
}

func NewSwaggerHandler(
	conf config.Swagger,
) *SwaggerHandler {
	return &SwaggerHandler{
		conf: conf,
	}
}

// RegisterSwagger registers swagger handler
func (h *SwaggerHandler) RegisterSwagger(e *echo.Echo) {
	docs.SwaggerInfo.Host = h.conf.Host
	docs.SwaggerInfo.Schemes = h.conf.Schemes
	e.GET(fmt.Sprintf("%s/*", h.conf.URL), echoSwagger.WrapHandler)
}
