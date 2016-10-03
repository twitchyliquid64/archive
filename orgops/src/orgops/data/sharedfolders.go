package data

import (
	"github.com/cznic/ql"
	"strconv"
)



var addSharedFolderQry = ql.MustCompile("BEGIN TRANSACTION; INSERT INTO sharedFolders(name,groupname) VALUES ($1, $2); COMMIT;")

func AddSharedFolder(name, groupname string)error{	//NOTE::This function should only be called from logic.Admin_NewSharedFolder, as this function
													//only sets up the DB for the shared folder.
	_, _, err := DB.Execute(ql.NewRWCtx(), addSharedFolderQry, name, groupname)
	return err
}


var getAllFoldersQry = ql.MustCompile("SELECT name,groupname FROM sharedFolders;")

func GetSharedFolders()(err error,folders []*[]string){
	results, _, err := DB.Execute(ql.NewRWCtx(), getAllFoldersQry)
	if err != nil{
		return err, nil
	}
	
	folderInfo, err := results[0].Rows(4096,0)//get up to 4096 shared folders
	if err != nil{
		return err, nil
	}
	if folderInfo == nil{
		return nil, nil
	}
	for _, row := range folderInfo{
		_, count := CountGroupMembers(row[1].(string))
		folder := &[]string{row[0].(string), row[1].(string), strconv.Itoa(count)}
		folders = append(folders, folder)
	}
	return
}


var delSharedFolderQry = ql.MustCompile("BEGIN TRANSACTION; DELETE FROM sharedFolders WHERE groupname == $1; COMMIT;")

func DelSharedFolder(groupname string)error{
	_, _, err := DB.Execute(ql.NewRWCtx(), delSharedFolderQry, groupname)
	return err
}


var UsrCanAccessSharedFolderQry = ql.MustCompile("SELECT groups.uid, users.username FROM (SELECT id() as ID, username FROM users) AS users, sharedFolders, groups WHERE sharedFolders.name == $1 && sharedFolders.groupname == groups.groupname && groups.uid == users.ID && users.username == $2")

func CanAccessSharedFolder(foldername, username string)(error, bool){
	results, _, err := DB.Execute(ql.NewRWCtx(), UsrCanAccessSharedFolderQry, foldername, username)
	if err != nil{
		return err, false
	}
	
	Info, err := results[0].Rows(4096,0)//get up to 4096 shared folders
	
	if len(Info) > 0{
		return nil, true
	}
	return nil, false
}


var UsrSharedFoldersList = ql.MustCompile("SELECT sharedFolders.name FROM sharedFolders, groups WHERE sharedFolders.groupname == groups.groupname && groups.uid == $1")

func GetUserFolders(uid int64)(err error,folders []string){
	results, _, err := DB.Execute(ql.NewRWCtx(), UsrSharedFoldersList, uid)
	if err != nil{
		return err, nil
	}
	
	folderInfo, err := results[0].Rows(4096,0)//get up to 4096 shared folders
	if err != nil{
		return err, nil
	}
	if folderInfo == nil{
		return nil, nil
	}
	for _, row := range folderInfo{
		folders = append(folders, row[0].(string))
	}
	return
}

