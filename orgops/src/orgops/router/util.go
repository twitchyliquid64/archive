package router

import (
	"github.com/hoisie/web"
	"orgops/templates"
	"orgops/session"
	"fmt"
)

func renderTemplate(ctx *web.Context, name string, data interface{}){
	t := templates.GetTemplate(name)
	if t==nil{ panic("No template found: "+name) }
	err := t.Execute(ctx.ResponseWriter, data)
	if err!=nil{
		fmt.Println(err)
	}
}


func getSessionKey(ctx *web.Context)string{
	for _, cookie := range ctx.Request.Cookies(){
		if cookie.Name == "sid"{
			return cookie.Value
		}
	}
	return ""
}

func deleteSessionKey(ctx *web.Context){
	for _, cookie := range ctx.Request.Cookies(){
		if cookie.Name == "sid"{
			cookie.MaxAge = -1
		}
	}
}

func sessionDetails(ctx *web.Context)map[string]string{
	return session.GetSession(getSessionKey(ctx))
}

func isAdmin(session map[string]string)bool{
	return session["admin"] == "yes"
}
