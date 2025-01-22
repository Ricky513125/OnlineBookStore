package services

import (
	"net/http"
	"strconv"

	"github.com/Bit0r/online-store/model"
	"github.com/Bit0r/online-store/model/perm"
	"github.com/Bit0r/online-store/view"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func ShowUsers(ctx *gin.Context) {
	isAdmin := ctx.Query("role") == "admin"
	step := uint64(10)

	paging := ctx.MustGet("paging").(view.Paging)
	count := model.CountUsers(isAdmin)
	paging.Total = count/step + 1
	ctx.Set("paging", paging)

	limit := model.Limit{Offset: (paging.Cur - 1) * step, Count: step}
	users, err := model.GetUsers(isAdmin, limit)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Set("tpl_files", []string{"layout.html", "navbar.html", "users.html"})
	ctx.Set("tpl_data", gin.H{
		"Users":   users,
		"HasPriv": lo.Contains[string],
		"IsAdmin": isAdmin,
	})
}

func EditPriv(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user, err := model.GetUser(userID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Set("tpl_files", []string{"layout.html", "navbar.html", "priv-editor.html"})
	ctx.Set("tpl_data", gin.H{
		"User":    user,
		"HasPriv": lo.Contains[string],
	})
}

func GrantPriv(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.PostForm("user_id"), 10, 64)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	privileges := ctx.PostFormArray("privileges")
	selfPrivileges := ctx.MustGet("privileges").(perm.PrivilegeSet)
	for _, priv := range privileges {
		// 检查自身是否有相应权限授权给其他用户
		if !selfPrivileges.HasPrivilege(priv) {
			ctx.Status(http.StatusForbidden)
			return
		}
	}

	err = model.SetPrivileges(userID, privileges)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "./user")
}
