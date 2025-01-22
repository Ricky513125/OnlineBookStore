package view

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Paging struct {
	Pre, Cur, Next, Total uint64
}

func Pagination(ctx *gin.Context) {
	paging := Paging{Total: 0}
	current, err := strconv.Atoi(ctx.Query("page"))
	if err == nil && current != 0 {
		paging.Cur = uint64(current)
	} else {
		paging.Cur = 1
	}
	ctx.Set("paging", paging)

	ctx.Next()

	paging = ctx.MustGet("paging").(Paging)
	if paging.Total == 0 {
		// 确认是否需要分页
		ctx.Set("paging", nil)
		return
	}

	// 设置当前页
	if paging.Cur > paging.Total {
		paging.Cur = paging.Total
	}
	paging.Pre, paging.Next = paging.Cur-1, paging.Cur+1
	ctx.Set("paging", paging)
}
