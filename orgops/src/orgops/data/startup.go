package data


import (
	"github.com/cznic/ql"
	"strings"
	"fmt"
)

var defaultIndexes map[string]string = map[string]string { //cant create index by type time
		"usersID:users": "id()",
		"usersUsername:users": "username",
		"usersAPI:users": "apikey",
		
		"groupsName:groups": "groupname",
		"groupsUid:groups": "uid",
		"groupsID:groups": "id()",
		
		"todoUid:todo": "uid",
		"todoID:todo": "id()",
		"todoCompleted:todo": "completed",
		
		"eventUid:event": "uid",
		"eventID:event": "id()",
		
		"feedUid:feeditems": "uid",
		"feedID:feeditems": "id()",
}


var defaultTables map[string][]string = map[string][]string {
		"users": []string {
			"username:string",
			"passhash:string",
			"otpsecret:string",
			"apikey:string",
		},
		"groups": []string {
			"uid:int64",
			"groupname:string",
		},
		"feeditems": []string { //These are like stories on a users' feed/dashboard.
			"uid:int64",
			"title:string",
			"content:string",
			"created:time",
		},
		"todo": []string { //Things a user has scheduled to complete.
			"uid:int64",
			"title:string",
			"description:string",
			"created:time",
			"completed:bool",
			"list:string",
		},
		"event": []string { //Things a user has scheduled to complete.
			"uid:int64",
			"title:string",
			"description:string",
			"created:time",
			"list:string",
			"eventtime:time",
		},//uid, title, description, list, created, eventtime
		"sharedFolders": []string {
			"name:string",
			"groupname:string",
		},
	}

//checks if all the default tables exist, and if they dont, it creates them.
func SetupDatabase()error {
	fmt.Println("(db) Now checking database for table integrity.")
	info, err := DB.Info()
	if err != nil{
		return err
	}
	
	//check the database has all the tables which are needed as a minimum
	for name, fields := range defaultTables {
		fmt.Println("(db) Checking: "+name)
		exists, _ := fetchTable(info.Tables , name)
		if !exists{
			err = createTable(DB, name, fields)
		}else{
		}
		
		if err != nil{
			return err
		}
	}
	
	
	//create the query to create the indexes if they dont exist
	fmt.Println("(db) Now checking/creating table indexes.")
	var sql string = "BEGIN TRANSACTION;\n"
	for name, field := range defaultIndexes {
		spl := strings.Split(name, ":")
		sql += "CREATE INDEX IF NOT EXISTS " + spl[0] + " ON " + spl[1] + " (" + field + ");\n"
	}
	sql += "COMMIT;"
	l, err := ql.Compile(sql)
	if err != nil{
		return err
	}
	_, _, err = DB.Execute(ql.NewRWCtx(), l)
	return err
}

func createTable(db *ql.DB, name string, fields []string) error{
	fmt.Println("     -Creating: "+name)
	var sql string = "BEGIN TRANSACTION;\n"
	sql += "CREATE TABLE " + name + "\n(\n"
	for _, field := range fields{
		spl := strings.Split(field, ":")
		sql += "  " + spl[0] + " " + spl[1] + ",\n"
	}
	sql += ");COMMIT;"
	l, err := ql.Compile(sql)
	if err != nil{
		return err
	}
	_, _, err = db.Execute(ql.NewRWCtx(), l)
	return err
}

//given the name of a table, returns true and the table if it exists in the table slice.
func fetchTable(tables []ql.TableInfo, name string)(bool,ql.TableInfo) {
	for _, table := range tables{
		if table.Name == name{
			return true, table
		}
	}
	return false, ql.TableInfo{}
}

