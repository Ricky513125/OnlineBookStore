package controller

import (
	"github.com/Ricky513125/OnlineBookStore/middleware"
	"github.com/Ricky513125/OnlineBookStore/services"
)

func setupAdmin() {
	adminGroup := router.Group("/admin", middleware.AuthUser)
	adminGroup.GET("/orders", middleware.Permission("order"), getOrders(true))
	adminGroup.GET("/user", middleware.Permission("user"), services.ShowUsers)
	adminGroup.GET("/user/:id", middleware.Permission("user"), services.EditPriv)
	adminGroup.POST("/user", middleware.Permission("user"), services.GrantPriv)
}
