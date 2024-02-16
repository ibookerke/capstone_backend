package http

import (
	"fmt"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
func (h *SwaggerHandler) RegisterSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Host = h.conf.Host
	docs.SwaggerInfo.Schemes = h.conf.Schemes
	r.GET(fmt.Sprintf("%s/*", h.conf.URL), ginSwagger.WrapHandler(swaggerfiles.Handler))
}
