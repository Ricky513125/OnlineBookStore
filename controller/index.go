package controller

import (
	"github.com/Bit0r/online-store/model"
	"github.com/Bit0r/online-store/model/perm"
	"github.com/Bit0r/online-store/view"
	"github.com/gin-gonic/gin"
)

func setupBooks() {
	router.GET("/index/*category", func(ctx *gin.Context) {
		var step uint64 = 12
		data := struct {
			IsAdmin    bool
			Category   string
			Categories []string
			model.Books
		}{}

		privileges, ok := ctx.Get("privileges")
		if ok && privileges.(perm.PrivilegeSet).HasPrivilege("book") {
			data.IsAdmin = true
		}

		// 填充过滤信息
		data.Category = ctx.Param("category")[1:]
		data.Categories = model.GetCategories()
		info := ctx.Query("info")

		// 填充过滤器
		filter := model.BooksFilter{
			Category:   data.Category,
			Info:       info,
			MustExists: !data.IsAdmin,
		}

		// 填充分页信息
		paging := ctx.MustGet("paging").(view.Paging)
		paging.Total = model.CountBooks(filter)/step + 1
		ctx.Set("paging", paging)

		// 填充图书信息
		data.Books, _ = model.GetBooks(
			filter,
			uint64((paging.Cur-1)*step),
			uint64(step))
		ctx.Set("tpl_data", data)

		// 设置模板
		files := []string{"layout.html", "home.html", "navbar.html"}
		ctx.Set("tpl_files", files)
	})
}
