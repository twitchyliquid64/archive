package router

import (
	"github.com/hoisie/web"
	"orgops/session"
	"orgops/logic"
	"fmt"
)



func Events_newevent(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "new_event.tmpl", DashboardPageData(s, nil, nil))
	}
}



func Events_neweventHandler(ctx *web.Context) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err := logic.NewEvent(session.GetUid(getSessionKey(ctx)), ctx.Params["title"], ctx.Params["desc"], ctx.Params["list"], ctx.Params["date"])
		if err!=nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/")
	}
}


func Events_viewevent(ctx *web.Context, eventid string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err, event := logic.SeeUsrEvent(session.GetUid(getSessionKey(ctx)), eventid)
		if err!=nil{
			fmt.Println(err)
		}
		renderTemplate(ctx, "view_event.tmpl", DashboardPageData(s, event, nil))
	}
}


func Events_deleteevent(ctx *web.Context, eventid string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		err := logic.DeleteEvent(session.GetUid(getSessionKey(ctx)), eventid)
		if err!=nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/")
	}
}
