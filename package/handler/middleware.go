package handler

import (
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

	if header != "" {
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
