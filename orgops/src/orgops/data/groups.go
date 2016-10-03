package data

import (
	"github.com/cznic/ql"
)




var getUsrGrpQry = ql.MustCompile("SELECT DISTINCT groupname FROM groups WHERE uid == $1;")

func GroupsByUserUid(uid int64)(error, []string){
	var groups []string
	results, _, err := DB.Execute(ql.NewRWCtx(), getUsrGrpQry, uid)
	if err != nil{
		return err, nil
	}
	
	groupData, err := results[0].Rows(2048,0)//max 2048 groups
	if err != nil{
		return err, nil
	}
	
	for _, group := range groupData {
		groups = append(groups, group[0].(string))
	}
	return nil, groups
}


var addUsrtoGrpQry = ql.MustCompile("BEGIN TRANSACTION; INSERT INTO groups(uid,groupname) VALUES ($1, $2); COMMIT;")

func AddUserToGroup(useruid int64, groupname string)error{
	_, _, err := DB.Execute(ql.NewRWCtx(), addUsrtoGrpQry, useruid, groupname)
	return err
}

func AddUserToGroups(useruid int64, groups []string)error{
	for _, group := range groups {
		err := AddUserToGroup(useruid, group)
		if err != nil{
			return err
		}
	}
	return nil
}

func AddUserToGroupByUsername(username string, groupname string)error{
	err, uid := UsernameToUid(username)
	if err!=nil{
		return err
	}
	return AddUserToGroup(uid, groupname)
}


var delUsrtoGrpQry = ql.MustCompile("BEGIN TRANSACTION; DELETE FROM groups WHERE uid == $1 && groupname == $2; COMMIT;")

func DelUserFromGroup(useruid int64, groupname string)error{
	_, _, err := DB.Execute(ql.NewRWCtx(), delUsrtoGrpQry, useruid, groupname)
	return err
}


func DelUserFromGroupByUsername(username string, groupname string)error{
	err, uid := UsernameToUid(username)
	if err!=nil{
		return err
	}
	return DelUserFromGroup(uid, groupname)
}


var countGrpMembersQry = ql.MustCompile("SELECT count() FROM groups WHERE groupname == $1;")

func CountGroupMembers(groupname string)(error, int){
	results, _, err := DB.Execute(ql.NewRWCtx(), countGrpMembersQry, groupname)
	if err != nil{
		return err, 0
	}
	
	Data, err := results[0].Rows(2048,0)//max 2048 groups
	if err != nil{
		return err, 0
	}
	
	
	return nil, int(Data[0][0].(int64))
}
