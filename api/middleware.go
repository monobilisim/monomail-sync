package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

func requireSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		store := ginsession.FromContext(ctx)
		if _, ok := store.Get("user"); !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		ctx.Next()
	}
}
