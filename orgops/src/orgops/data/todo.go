package data

import (
	"github.com/cznic/ql"
	"time"
)


type Todo struct {
	UID int64
	Title string
	Description string
	Completed bool
	Created time.Time
	CreatedBy int64
	List string
}


var getTodoQry = ql.MustCompile("SELECT title, description, created, list, completed, uid FROM todo WHERE id() == $1;")
func GetTodo(uid int64)(error, *Todo){
	results, _, err := DB.Execute(ql.NewRWCtx(), getTodoQry, uid)
	if err != nil{
		return err, nil
	}
	
	todo, err := results[0].FirstRow()
	if err != nil{
		return err, nil
	}
	
	if todo == nil{//no result
		return nil,nil
	}
	
	return nil, &Todo{
			UID: uid,
			Title: todo[0].(string),
			Description: todo[1].(string),
			Completed: todo[4].(bool),
			Created: todo[2].(time.Time),
			List: todo[3].(string),
			CreatedBy: todo[5].(int64),
	}
}


var markCompleteTodoQry = ql.MustCompile("BEGIN TRANSACTION; UPDATE todo SET completed=true WHERE id() == $1 && uid == $2; COMMIT;")

func (inst *Todo)Complete(owneruid int64)error{
	_, _, err := DB.Execute(ql.NewRWCtx(), markCompleteTodoQry, inst.UID, owneruid)
	return err
}


var getUsrTodoQry = ql.MustCompile("SELECT id(), title, description, created, list FROM todo WHERE uid == $1 && completed == false ORDER BY created DESC LIMIT 20;")

func TodoByUserUid(uid int64)(error, []Todo){
	var ret []Todo
	results, _, err := DB.Execute(ql.NewRWCtx(), getUsrTodoQry, uid)
	if err != nil{
		return err, nil
	}
	
	todoData, err := results[0].Rows(2048,0)//max 2048 todo's
	if err != nil{
		return err, nil
	}
	
	for _, todo := range todoData {
		ret = append(ret, Todo{
			UID: todo[0].(int64),
			Title: todo[1].(string),
			Description: todo[2].(string),
			Completed: false,
			Created: todo[3].(time.Time),
			List: todo[4].(string),
			CreatedBy: uid,//users uid, not task uid
		})
	}
	return nil, ret
}


var addTodoQry = ql.MustCompile("BEGIN TRANSACTION; INSERT INTO todo(uid,title,description,created,completed,list) VALUES ($1, $2, $3, now(), false, $4); COMMIT;")

func AddTodo(useruid int64, title, description, list string)error{
	_, _, err := DB.Execute(ql.NewRWCtx(), addTodoQry, useruid, title, description, list)
	return err
}
