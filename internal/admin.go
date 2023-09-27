package internal

import "github.com/gin-gonic/gin"

func handleAdmin(ctx *gin.Context) {
	ctx.HTML(200, "admin.html", nil)
}
