package router

import (
	"github.com/hoisie/web"
	"orgops/logic"
	"orgops/data"
	"fmt"
)


func Admin(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		renderTemplate(ctx, "admin.tmpl", AdminPageData(s, logic.Status(), nil))
	}
}

func Admin_users(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err, obj := logic.GetUserList()
		if err!=nil{
			fmt.Println(err)
		}
		renderTemplate(ctx, "useradmin.tmpl", AdminPageData(s, obj, nil))
	}
}

func Admin_viewUser(ctx *web.Context, user string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err, groups := logic.UserGroups(user)
		if err != nil{
			fmt.Println(err)
		}
		renderTemplate(ctx, "useradmin_edit.tmpl", AdminPageData(s, UserUpdateData{Username: user, IsAdmin: logic.IsAdmin(user), Groups: groups }, nil))
	}
}


func Admin_updateUserHandler(ctx *web.Context, user string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err := logic.Admin_UpdateUser(user, ctx.Params["addgroup"], ctx.Params["removegroup"])
		err2, groups := logic.UserGroups(user)
		if err2 != nil{
			fmt.Println(err)
		}
		renderTemplate(ctx, "useradmin_edit.tmpl", AdminPageData(s, UserUpdateData{Username: user, IsAdmin: logic.IsAdmin(user), Groups: groups }, err))
	}
}

func Admin_changepw(ctx *web.Context, user string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{

		renderTemplate(ctx, "useradmin_changepw.tmpl", AdminPageData(s, UserUpdateData{Username: user, }, nil))
	}
}

func Admin_changepwHandler(ctx *web.Context, user string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err := logic.Admin_ChangePw(user, ctx.Params["passwd"])
		if err==nil{
			ctx.Redirect(302, "/admin/users")
		}else{
			renderTemplate(ctx, "useradmin_changepw.tmpl", AdminPageData(s, UserUpdateData{Username: user, }, err))
		}
	}
}


func Admin_newuser(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		renderTemplate(ctx, "useradmin_newuser.tmpl", AdminPageData(s, nil, nil))
	}
}


func Admin_datastorage(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err, headings, rows := logic.Admin_DoQuery(ctx.Params["query"])
		err2, folders := data.GetSharedFolders()
		if err2 != nil{
			fmt.Println(err2)
		}
		renderTemplate(ctx, "admin_database.tmpl", AdminPageData(s, DataStoragePageData(ctx.Params["query"], headings, rows, folders), err))
	}
}

func Admin_newuserhandler(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err := logic.Admin_NewUser(ctx.Params["usr"], ctx.Params["passwd"])
		if err==nil{
			ctx.Redirect(302, "/admin/users")
		}else{
			renderTemplate(ctx, "useradmin_newuser.tmpl", AdminPageData(s, err, err))
		}
	}
}



func Admin_newsharedfolderhandler(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err := logic.Admin_NewSharedFolder(ctx.Params["name"], ctx.Params["gname"])
		if err==nil{
			ctx.Redirect(302, "/admin/repos")
		}else{
			err2, folders := data.GetSharedFolders()
			if err2 != nil{
				fmt.Println(err2)
			}
			renderTemplate(ctx, "admin_database.tmpl", AdminPageData(s, DataStoragePageData("", nil, nil, folders), err))
		}
	}
}


func Admin_deletesharedfolderhandler(ctx *web.Context, groupname string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else if !isAdmin(s){
		ctx.Redirect(302, "/")
	}else{
		err := logic.Admin_DeleteSharedFolder(groupname)
		err2, folders := data.GetSharedFolders()
		if err2 != nil{
			fmt.Println(err2)
		}
		renderTemplate(ctx, "admin_database.tmpl", AdminPageData(s, DataStoragePageData("", nil, nil, folders), err))
	}
}
