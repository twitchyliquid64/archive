package router

import (
	"github.com/hoisie/web"
	"encoding/json"
	"orgops/data"
	"errors"
)



func API_userdetails(ctx *web.Context, username, key string) {
	
	err, usr := getUserAPIData(ctx, username, key)
	if err != nil{
		return
	}
	usr.Passhash = "REDACTED"
	
	writeJSON(ctx, map[string]interface{}{"err": err, "user": usr})
}



func API_userupcoming(ctx *web.Context, username, key string) {
	
	err, usr := getUserAPIData(ctx, username, key)
	if err != nil{
		return
	}

	err, todo := data.TodoByUserUid(usr.UID)
	if err!=nil{
		data, _ := json.Marshal(map[string]interface{}{"err": err})
		ctx.Write(data)
		return
	}


	err, events := data.EventByUserUid(usr.UID)
	if err!=nil{
		data, _ := json.Marshal(map[string]interface{}{"err": err})
		ctx.Write(data)
		return
	}

	writeJSON(ctx, map[string]interface{}{"err": err, "todo": todo, "events": events})
}


func writeJSON(ctx *web.Context, in interface{}){
	data, err := json.Marshal(in)
	if err==nil{
		ctx.Write(data)
	}else{
		data, err = json.Marshal(map[string]interface{}{"err": err})
		ctx.Write(data)
	}
}


func getUserAPIData(ctx *web.Context, username string, key string)(error, *data.User){
	err, usr := data.GetUserByName(username)
	if err!=nil{
		data, _ := json.Marshal(map[string]interface{}{"err": err})
		ctx.Write(data)
		return err, nil
	}
	
	if usr.APIKey != key{
		data, _ := json.Marshal(map[string]interface{}{"err": "API Key mismatch"})
		ctx.Write(data)
		return errors.New("Key mismatch"), nil
	}
	return nil, usr
}
