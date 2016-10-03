package data

import (
	"github.com/cznic/ql"
)

type User struct{
	UID				int64
	Username 		string
	Passhash		string
	OtpSecret 		string
	APIKey			string
	Groups 			[]string
}



var getUsrQry = ql.MustCompile("SELECT DISTINCT id(), username, passhash, otpsecret, apikey FROM users WHERE username == $1;")

func GetUserByName(username string)(error,*User){
	results, _, err := DB.Execute(ql.NewRWCtx(), getUsrQry, username)
	if err != nil{
		return err, nil
	}
	
	usrInfo, err := results[0].FirstRow()
	if err != nil{
		return err, nil
	}
	
	if usrInfo == nil{//no result
		return nil,nil
	}
	
	err, groups := GroupsByUserUid(usrInfo[0].(int64))
	if err != nil{
		return err, nil
	}
	
	return nil, &User{
		UID:		usrInfo[0].(int64),
		Username: 	usrInfo[1].(string),
		Passhash: 	usrInfo[2].(string),
		OtpSecret: 	usrInfo[3].(string),
		APIKey:		usrInfo[4].(string),
		Groups:		groups,
	}
}


var getUsrByUidQry = ql.MustCompile("SELECT username FROM users WHERE id() == $1;")

func GetUserByUid(uid int64)(error,*User){
	results, _, err := DB.Execute(ql.NewRWCtx(), getUsrByUidQry, uid)
	if err != nil{
		return err, nil
	}
	
	usrInfo, err := results[0].FirstRow()
	if err != nil{
		return err, nil
	}
	
	if usrInfo == nil{//no result
		return nil,nil
	}
	return GetUserByName(usrInfo[0].(string))
}


var createUsrQry = ql.MustCompile("BEGIN TRANSACTION; INSERT INTO users(username,passhash,otpsecret,apikey) VALUES ($1, $2, $3, $4); COMMIT;")

func NewUser(usr *User)error{
	
	if usr.APIKey == "" {
		usr.APIKey = randString(16)
	}
	
	context := ql.NewRWCtx()
	_, _, err := DB.Execute(context, createUsrQry, usr.Username, usr.Passhash, usr.OtpSecret, usr.APIKey)
	if err != nil{
		return err
	}
	
	uid := context.LastInsertID
	usr.UID = uid
	err = AddUserToGroups(uid, usr.Groups)
	return err
}


var getUidByUserQry = ql.MustCompile("SELECT id() FROM users WHERE username == $1;")

func UsernameToUid(username string)(error, int64){
	results, _, err := DB.Execute(ql.NewRWCtx(), getUidByUserQry, username)
	if err != nil{
		return err, 0
	}
	
	usrInfo, err := results[0].FirstRow()
	if err != nil{
		return err, 0
	}
	
	if usrInfo == nil{//no result
		return nil,0
	}
	return nil, usrInfo[0].(int64)
}

var commitUsrQry = ql.MustCompile("BEGIN TRANSACTION; UPDATE users SET username=$1,passhash=$2,otpsecret=$3,apikey=$4 WHERE id() == $5 ; COMMIT;")

func (usr *User)Commit()error{
	_, _, err := DB.Execute(ql.NewRWCtx(), commitUsrQry, usr.Username, usr.Passhash, usr.OtpSecret, usr.APIKey, usr.UID)
	if err != nil{
		return err
	}
	return nil
}


var getAllUsrsQry = ql.MustCompile("SELECT id() FROM users;")

func GetUsers()(err error,users []*User){
	results, _, err := DB.Execute(ql.NewRWCtx(), getAllUsrsQry)
	if err != nil{
		return err, nil
	}
	
	usrInfo, err := results[0].Rows(4096,0)//get up to 4096 users
	if err != nil{
		return err, nil
	}
	if usrInfo == nil{
		return nil, nil
	}
	for _, user := range usrInfo{
		err, usr := GetUserByUid(user[0].(int64))
		if err != nil{
			return err, nil
		}
		users = append(users, usr)
	}
	return
}

var UsrCountQry = ql.MustCompile("SELECT count() FROM users;")

func UsrCount()int{
	results, _, err := DB.Execute(ql.NewRWCtx(), getAllUsrsQry)
	if err != nil{
		return 0
	}
	
	usrInfo, err := results[0].FirstRow()//get up to 4096 users
	if err != nil{
		return 0
	}
	if usrInfo == nil{
		return 0
	}
	return int(usrInfo[0].(int64))
}
