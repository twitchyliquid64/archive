package router

import (
	"github.com/hoisie/web"
	"orgops/logic"
)


func Install(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		renderTemplate(ctx, "install.tmpl", nil)
	}else{
		ctx.Redirect(302, "/login")
	}
}

func InstallHandler(ctx *web.Context) {

	ok := logic.SetupAdminAccount(ctx.Params["username"],ctx.Params["passwd"])
	if ok{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "install.tmpl", "Unspecified error.")
	}
}

