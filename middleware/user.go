package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthUser(ctx *gin.Context) {
	if !isLogged(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func AuthUserRedirect(ctx *gin.Context) {
	if !isLogged(ctx) {
		ctx.Redirect(http.StatusFound, "/user/log-in")
		ctx.Abort()
	}
}

func isLogged(ctx *gin.Context) bool {
	return ctx.GetBool("isLoggedIn")
}

func LogIn(ctx *gin.Context, userID uint64) {
	session := sessions.Default(ctx)
	session.Set("userID", userID)
	session.Save()
}

func LogOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err == nil {
		ctx.Redirect(http.StatusFound, "/")
	} else {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
	}
}
