package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Bit0r/online-store/middleware"
	"github.com/Bit0r/online-store/model"
	"github.com/gin-gonic/gin"
)

func setupCart() {
	cartGroup := router.Group("/cart")
	cartGroup.POST("/book", middleware.AuthUser, handleAddCart)
	cartGroup.GET("/books", middleware.AuthUserRedirect, handleShowCart)
	cartGroup.PUT("/book", middleware.AuthUser, handleUpdateCart)
}

func handleShowCart(ctx *gin.Context) {
	id := ctx.GetUint64("userID")

	files := []string{"layout.html", "navbar.html", "shopping-cart.html"}

	addresses, err := model.GetAddresses(id)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	data := struct {
		model.CartBooks
		model.Addresses
	}{
		model.GetCartItems(id),
		addresses}

	ctx.Set("tpl_files", files)
	ctx.Set("tpl_data", data)
}

func handleAddCart(ctx *gin.Context) {
	id := ctx.GetUint64("userID")

	bookID, _ := strconv.ParseUint(ctx.PostForm("id"), 0, 64)
	err := model.AddCartItem(id, bookID)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
	}
}

func handleUpdateCart(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	bookInfo := struct {
		BookID   uint64 `form:"id"`
		Quantity uint   `form:"quantity"`
	}{}
	ctx.Bind(&bookInfo)
	err := model.UpdateCartItem(userID, bookInfo.BookID, bookInfo.Quantity)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
	}
}
