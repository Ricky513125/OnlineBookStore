package middleware

import (
	"log"
	"net/http"

	"github.com/Bit0r/online-store/model"
	"github.com/Bit0r/online-store/model/perm"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetupPrivileges(ctx *gin.Context) {
	userID, ok := sessions.Default(ctx).Get("userID").(uint64)
	ctx.Set("isLoggedIn", ok)
	if !ok {
		return
	}
	ctx.Set("userID", userID)

	privileges, err := model.GetPrivilegeSet(userID)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.Set("privileges", privileges)
}

func Permission(privilege string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		privileges, ok := ctx.Get("privileges")
		if !ok {
			ctx.AbortWithStatus(http.StatusForbidden)
		}

		if !privileges.(perm.PrivilegeSet).HasPrivilege(privilege) {
			ctx.AbortWithStatus(http.StatusForbidden)
		}
	}
}
