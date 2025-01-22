package controller

import (
	"net/http"

	"github.com/Bit0r/online-store/middleware"
	"github.com/Bit0r/online-store/model"
	"github.com/Bit0r/online-store/services"
	"github.com/gin-gonic/gin"
)

func setupUser() {
	userGroup := router.Group("/user")
	userGroup.Any("/log-in", handleLogin)
	userGroup.Any("/sign-up", handleSignUp)
	userGroup.GET("/orders", middleware.AuthUserRedirect, getOrders(false))
	userGroup.GET("/log-out", middleware.AuthUser, handleLogout)

	addressGroup := userGroup.Group("/address")
	addressGroup.GET("/", middleware.AuthUserRedirect, services.GetAddresses)
	addressGroup.POST("/", middleware.AuthUserRedirect, services.ReplaceAddress)
	addressGroup.GET("/:id", middleware.AuthUserRedirect, services.EditAddress)
	addressGroup.DELETE("/:id", middleware.AuthUser, services.DeleteAddress)
}

func handleLogin(ctx *gin.Context) {
	switch ctx.Request.Method {
	case http.MethodGet:
		files := []string{"layout-no-nav.html", "log-in.html"}
		ctx.Set("tpl_files", files)
	case http.MethodPost:
		id := model.VerifyUser(ctx.PostForm("name"), ctx.PostForm("passwd"))
		if id != 0 {
			middleware.LogIn(ctx, id)
		} else {
			ctx.Status(http.StatusUnauthorized)
		}
	}
}

func handleSignUp(ctx *gin.Context) {
	switch ctx.Request.Method {
	case http.MethodGet:
		files := []string{"layout-no-nav.html", "sign-up.html"}
		ctx.Set("tpl_files", files)
	case http.MethodPost:
		id := model.AddUser(ctx.PostForm("name"), ctx.PostForm("passwd"))
		if id != 0 {
			middleware.LogIn(ctx, id)
		} else {
			ctx.Status(http.StatusUnauthorized)
		}
	}
}

func handleLogout(ctx *gin.Context) {
	middleware.LogOut(ctx)
}
