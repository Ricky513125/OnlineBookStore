package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Bit0r/online-store/middleware"
	"github.com/Bit0r/online-store/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupBook() {
	bookGroup := router.Group("/book")

	bookGroup.GET("/edit/*id",
		middleware.AuthUserRedirect,
		middleware.Permission("book"),
		bookEditor)

	bookGroup.POST("/edit",
		middleware.AuthUserRedirect,
		middleware.Permission("book"),
		bookUpdate)
}

func bookUpdate(ctx *gin.Context) {
	bookInfo := struct {
		ID         uint64  `form:"id" binding:"required"`
		ISBN       string  `form:"isbn" binding:"required"`
		Name       string  `form:"name" binding:"required"`
		Author     string  `form:"author" binding:"required"`
		Price      float64 `form:"price" binding:"required"`
		Intro      string  `form:"intro"`
		Deleted    string  `form:"deleted"`
		Categories string  `form:"categories"`
	}{}

	ctx.Bind(&bookInfo)

	book := model.Book{
		ID:      bookInfo.ID,
		ISBN:    bookInfo.ISBN,
		Name:    bookInfo.Name,
		Author:  bookInfo.Author,
		Price:   bookInfo.Price,
		Deleted: bookInfo.Deleted != "",
	}

	if bookInfo.Intro != "" {
		// 如果有intro，则设置intro
		book.Intro = sql.NullString{
			String: bookInfo.Intro,
			Valid:  true,
		}
	}

	// 如果有传入cover，则更新cover
	if file, err := ctx.FormFile("cover"); err == nil {
		// 检验文件类型
		mimeType := file.Header.Get("Content-Type")
		if mimeType[:5] != "image" {
			// 如果上传的不是图片，则返回错误
			ctx.Status(http.StatusBadRequest)
			return
		}

		// 保存封面图片
		fileName := fmt.Sprintf("%s.%s", uuid.New(), mimeType[6:])
		dst := fmt.Sprintf("%s/cover/%s", uploadDir, fileName)
		err = ctx.SaveUploadedFile(file, dst)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if book.ID != 0 {
			// 如果是更新，则删除旧的封面
			if oldBook, _ := model.GetBook(book.ID); oldBook.Cover.Valid {
				// 如果之前有封面，则删除之前的封面
				os.Remove(fmt.Sprintf("%s/cover/%s", uploadDir, oldBook.Cover.String))
			}
		}

		// 设置封面
		book.Cover = sql.NullString{
			String: fileName,
			Valid:  true,
		}
	}

	if book.ID == 0 {
		// 如果id为0，则是新增
		book.ID, _ = model.AddBook(book)
	} else {
		// 如果id不为0，则是更新
		model.UpdateBook(book)
	}

	if bookInfo.Categories != "" {
		model.UpdateCategories(book.ID, strings.Split(bookInfo.Categories, ","))
	}

	ctx.Redirect(http.StatusFound, "./edit/"+strconv.FormatUint(book.ID, 10))
}

func bookEditor(ctx *gin.Context) {
	ctx.Set("tpl_files", []string{"layout.html", "navbar.html", "book-editor.html"})

	bookID, err := strconv.ParseUint(ctx.Param("id")[1:], 10, 64)
	if err != nil {
		// 如果没有传入id，则填充一个空的book
		ctx.Set("tpl_data", struct {
			model.Book
			Categories string
		}{})
		return
	}

	// 如果传入了id，则填充book
	book, err := model.GetBook(bookID)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Set("tpl_data", struct {
		model.Book
		Categories string
	}{book, strings.Join(model.GetBookCategories(bookID), ",")})
}
