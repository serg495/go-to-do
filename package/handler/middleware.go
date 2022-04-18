package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authHeader = "Authorization"
	userCtx    = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		newErrorMessage(c, http.StatusUnauthorized, "missing auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorMessage(c, http.StatusUnauthorized, "auth header is invalid")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorMessage(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorMessage(c, http.StatusNotFound, "user id not found")
		return 0, errors.New("user id not found")
	}

	intId, ok := id.(int)
	if !ok {
		newErrorMessage(c, http.StatusInternalServerError, "format id is invalid")
		return 0, errors.New("format id is invalid")
	}

	return intId, nil

}
