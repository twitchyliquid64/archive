package logic

import (
	"orgops/session"
    "encoding/base64"
	"crypto/sha1"
	"orgops/data"
	"errors"
)

func Login(user, pass, otpkey string)(err error, sessionkey string){
	err, usr := data.GetUserByName(user)
	if err!=nil{
		return errors.New("Database Error: "+err.Error()), ""
	}
	
	ok := checkPass(usr, pass)
	
	if !ok{//login failure
		return errors.New("Authentication failure. Please check your username/password combination."), ""
	}
	
	adminstr := "no"
	if chkStrInStrSlice("admin", usr.Groups){
		adminstr = "yes"
	}
	
	sessionkey = session.CreateSession(usr.UID, map[string]string{
					"username": user,
					"admin": adminstr,
				})
	return
}



func ChangePassword( user, oldpass, newpass, retypepass string, override bool)error{
	if newpass != retypepass{
		return errors.New("The password you retyped did not match your new password.")
	}
	err, usr := data.GetUserByName(user)
	if err!=nil{
		return errors.New("Database Error: "+err.Error())
	}
	
	if !override{
		if shaHash(user, oldpass) != usr.Passhash {
			return errors.New("Incorrect password.")
		}
	}
	
	usr.Passhash = shaHash(user, newpass)
	return usr.Commit()
}

func UserGroups(user string)(error, []string){
	err, usr := data.GetUserByName(user)
	if err!=nil{
		return errors.New("Database Error: "+err.Error()), nil
	}
	if usr == nil{return nil, nil}
	return nil, usr.Groups
}

func UserAPIKey(user string)(error, string){
	err, usr := data.GetUserByName(user)
	if err!=nil{
		return errors.New("Database Error: "+err.Error()), ""
	}
	if usr == nil{return nil, ""}
	return nil, usr.APIKey
}


func Logout(key string){
	session.DeleteSession(key)
}




func checkPass(usr *data.User, pass string)bool{
    return shaHash(usr.Username, pass) == usr.Passhash
}

func shaHash(usr, pass string)string{
	hasher := sha1.New()
    hasher.Write([]byte(pass))
    hasher.Write([]byte(usr))
    return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
