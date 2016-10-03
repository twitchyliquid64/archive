package router

import (
	"github.com/hoisie/web"
	"orgops/session"
	"orgops/logic"
	"orgops/data"
	"fmt"
)

func Login(ctx *web.Context) {
	renderTemplate(ctx, "login.tmpl", nil)
}

func LoginHandler(ctx *web.Context) {
	err, skey := logic.Login(ctx.Params["username"], ctx.Params["passwd"], ctx.Params["otp"])
	
	if err != nil{
		renderTemplate(ctx, "login.tmpl", err.Error())
	}else{
		ctx.SetCookie(web.NewCookie("sid", skey, 60*60*24*5))
		ctx.Redirect(302, "/user/setting/account")
	}
}



func Logout(ctx *web.Context) {
	logic.Logout(getSessionKey(ctx))
	deleteSessionKey(ctx)
	renderTemplate(ctx, "login.tmpl", nil)
}


func Index(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err, todos := data.TodoByUserUid(session.GetUid(getSessionKey(ctx)))
		if err!=nil{
			fmt.Println(err)
		}
		err, events := data.EventByUserUid(session.GetUid(getSessionKey(ctx)))
		if err!=nil{
			fmt.Println(err)
		}
		renderTemplate(ctx, "dashboard.tmpl", DashboardPageData(s, todos, events))
	}
}


func Todo_newtodo(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "new_todo.tmpl", DashboardPageData(s, nil, nil))
	}
}


func Todo_viewtodo(ctx *web.Context, taskid string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err, todo := logic.GetUsrTodo(session.GetUid(getSessionKey(ctx)), taskid)
		if err!=nil{
			fmt.Println(err)
		}
		renderTemplate(ctx, "view_todo.tmpl", DashboardPageData(s, todo, nil))
	}
}


func Todo_newtodoHandler(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err := logic.NewTodo(session.GetUid(getSessionKey(ctx)), ctx.Params["title"], ctx.Params["desc"], ctx.Params["list"])
		if err!=nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/")
	}
}

func Todo_setcompleteHandler(ctx *web.Context, taskid string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err := logic.SetComplete(session.GetUid(getSessionKey(ctx)), taskid)
		if err!=nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/")
	}
}


func Settings_acctdetails(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err, groups := logic.UserGroups(s["username"])
		if err != nil{
			renderTemplate(ctx, "account.tmpl", UserSettingsPageData(s, err, "", groups, ""))
		}else{
			err, apikey := logic.UserAPIKey(s["username"])
			renderTemplate(ctx, "account.tmpl", UserSettingsPageData(s, err, "", groups, apikey))
		}
	}
}

func Settings(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "password.tmpl", UserSettingsPageData(s, nil, "", nil, ""))
	}
}


func ChangePasswordHandler(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err := logic.ChangePassword(s["username"], ctx.Params["oldpasswd"], ctx.Params["newpasswd"], ctx.Params["retypepasswd"], false)
		var successMsg string
		if err == nil{
			successMsg = "Password changed successfully."
		}
		renderTemplate(ctx, "password.tmpl", UserSettingsPageData(s, err, successMsg, nil, ""))
	}
}


func Settings_2f(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "2factor.tmpl", UserSettingsPageData(s, nil, "", nil, ""))
	}
}

func Settings_delete(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "delete.tmpl", UserSettingsPageData(s, nil, "", nil, ""))
	}
}
