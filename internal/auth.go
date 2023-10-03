package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

func handleLogin(ctx *gin.Context) {
	store := ginsession.FromContext(ctx)
	_, ok := store.Get("user")
	if ok {
		// User is not logged in, redirect them to the login page
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	ctx.HTML(200, "login.html", nil)
}

func login(ctx *gin.Context) {
	var data user
	data.Name = ctx.PostForm("username")
	data.Password = ctx.PostForm("password")

	pass, err := getPassword(data.Name)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(data.Password))
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "wrong password"})
		return
	}
	store := ginsession.FromContext(ctx)

	store.Set("user", data.Name)
	if err := store.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, "/admin")

}
