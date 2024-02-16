package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/ibookerke/capstone_backend/internal/domain"
)

func GetUserId(ctx *gin.Context) (domain.UserID, error) {
	claimsValue, exists := ctx.Get("claims")
	if !exists {
		return 0, errors.New("unauthorized")
	}

	claims, ok := claimsValue.(*Claims)
	if !ok {
		return 0, errors.New("unauthorized")
	}

	userID := claims.User.ID
	if userID == 0 {
		return 0, errors.New("unauthorized")
	}

	return userID, nil
}
