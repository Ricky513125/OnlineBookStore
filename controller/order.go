package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Bit0r/online-store/middleware"
	"github.com/Bit0r/online-store/model"
	"github.com/Bit0r/online-store/model/perm"
	"github.com/Bit0r/online-store/view"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var statusI18n = map[string]string{
	"unpaid":     "待付款",
	"failed":     "失败",
	"to_be_ship": "待发货",
	"shipped":    "已发货，正在运输",
	"success":    "成功",
}

func setupOrder() {
	router.GET("/order/:id", middleware.AuthUserRedirect, handleGetOrder)
	router.POST("/order", middleware.AuthUser, handleAddOrder)
	router.PUT("/order", middleware.AuthUser, handleUpdateOrder)
}

func getOrders(needAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		isAdmin := ctx.MustGet("privileges").(perm.PrivilegeSet).HasPrivilege("order")

		var orders model.Orders
		var userID uint64
		switch {
		case needAdmin && isAdmin:
			userID = 0
		case !needAdmin:
			userID = ctx.GetUint64("userID")
		default:
			ctx.Status(http.StatusForbidden)
		}

		// 获取分页信息
		step := uint64(3)
		paging := ctx.MustGet("paging").(view.Paging)
		offset := (paging.Cur - 1) * step

		// 获取订单信息
		orders, err := model.GetOrders(userID, offset, step)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// 获取订单数量
		count, err := model.CountOrders(userID)
		if err != nil {
			log.Println(err)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		// 设置分页信息
		paging.Total = count/step + 1
		ctx.Set("paging", paging)

		for idx := range orders {
			orders[idx].Status = statusI18n[orders[idx].Status]
		}

		ctx.Set("tpl_files", []string{"layout.html", "navbar.html", "orders.html"})
		ctx.Set("tpl_data", orders)
		if isAdmin && !needAdmin {
			ctx.Set("buttons",
				append(ctx.MustGet("buttons").(view.Buttons),
					view.Button{Text: "所有订单", URL: "/admin/orders", Class: "is-danger"}))
		}
	}
}

func handleGetOrder(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")
	orderID, _ := strconv.ParseUint(ctx.Param("id"), 0, 64)

	// 获取订单信息
	order, err := model.GetOrder(orderID, true)
	if err != nil {
		log.Println(err)
		return
	}

	// 校验访问权限
	isAdmin := ctx.MustGet("privileges").(perm.PrivilegeSet).HasPrivilege("order")
	if order.UserID != userID && !isAdmin {
		ctx.Status(http.StatusForbidden)
		return
	}

	// 获取地址信息
	address, _ := model.GetAddress(order.AddressID)

	// 访问模板
	ctx.Set("tpl_files", []string{"layout.html", "order.html", "navbar.html"})

	ctx.Set("tpl_data", struct {
		model.Order
		model.Address
		IsAdmin  bool
		StatusCN string
	}{order, address, isAdmin, statusI18n[order.Status]})
}

func handleAddOrder(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")

	addressID, _ := strconv.Atoi(ctx.PostForm("addressID"))

	booksID := []uint64{}
	for _, v := range ctx.PostFormArray("booksID") {
		bookID, _ := strconv.ParseUint(v, 0, 64)
		booksID = append(booksID, bookID)
	}

	orderID, err := model.AddOrder(userID, uint(addressID), booksID)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/order/"+strconv.Itoa(int(orderID)))
	}
}

func handleUpdateOrder(ctx *gin.Context) {
	userID := ctx.GetUint64("userID")

	var orderForm struct {
		ID     uint64 `form:"id"`
		Status string `form:"status"`
	}
	ctx.Bind(&orderForm)
	log.Println(orderForm.ID, orderForm.Status)
	// 获取订单信息
	order, err := model.GetOrder(orderForm.ID, false)
	if err != nil {
		log.Println(err)
		return
	}

	// 校验访问权限
	isAdmin := ctx.MustGet("privileges").(perm.PrivilegeSet).HasPrivilege("order")
	if order.UserID != userID && !isAdmin {
		ctx.Status(http.StatusForbidden)
		return
	}

	// 校验状态转移的正确性
	type edge struct {
		// cur       string
		next      string
		needAdmin bool
	}
	dag := map[string][]edge{
		"unpaid":     {{next: "to_be_ship"}, {next: "failed"}},
		"to_be_ship": {{next: "failed", needAdmin: true}, {next: "shipped", needAdmin: true}},
		"shipped":    {{next: "failed"}, {next: "success"}},
		"success":    {},
		"failed":     {},
	}
	success := lo.ContainsBy(dag[order.Status], func(e edge) bool {
		return e.next == orderForm.Status &&
			(!e.needAdmin || (e.needAdmin && isAdmin))
	})
	if !success {
		ctx.Status(http.StatusConflict)
	}

	// 更新订单状态
	err = model.UpdateOrderStatus(orderForm.ID, orderForm.Status)
	if err != nil {
		log.Println(err)
		ctx.Status(http.StatusInternalServerError)
	}
}
