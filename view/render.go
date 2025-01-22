package view

import (
	"html/template"
	"log"

	"github.com/Bit0r/online-store/model/perm"
	"github.com/gin-gonic/gin"
)

func TemplateExecute(ctx *gin.Context) {
	ctx.Set("buttons", Buttons{})

	ctx.Next()

	_, ok := ctx.Get("tpl_files")
	if !ok {
		return
	}

	files := ctx.GetStringSlice("tpl_files")
	files = append(files, "pagination.html")
	for idx, file := range files {
		files[idx] = tplRoot + file
	}
	hasPrivilege := func(privilege string) bool {
		return ctx.MustGet("privileges").(perm.PrivilegeSet).HasPrivilege(privilege)
	}

	tpl, _ := template.New("layout").
		Funcs(template.FuncMap{"hasPrivilege": hasPrivilege}).ParseFiles(files...)

	data := gin.H{}
	tpl_data, ok := ctx.Get("tpl_data")
	if !ok {
		err := tpl.Execute(ctx.Writer, nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	data["IsLoggedIn"] = ctx.GetBool("isLoggedIn")
	data["Data"] = tpl_data
	data["Paging"] = ctx.MustGet("paging")
	data["Buttons"] = ctx.MustGet("buttons")

	err := tpl.Execute(ctx.Writer, data)
	if err != nil {
		log.Println(err)
	}
}
