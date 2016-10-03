package logic

import(
	"orgops/data"
	"strconv"
	"errors"
)

func NewTodo(useruid int64, title, desc, list string)error{
	return data.AddTodo(useruid, title, desc, list)
}

func SetComplete(useruid int64, todouid string)error{
	
	tid, err := strconv.ParseInt(todouid, 10, 64)
	if err!=nil{
		return err
	}
	
	err, todo := data.GetTodo(tid)
	if err!=nil{
		return err
	}
	
	err = todo.Complete(useruid)//The SQL makes sure a user cannot set another users task as complete
	return err
}


func GetUsrTodo(useruid int64, todouid string)(error, *data.Todo){
	
	tid, err := strconv.ParseInt(todouid, 10, 64)
	if err!=nil{
		return err, nil
	}
	
	err, todo := data.GetTodo(tid)
	if err!=nil{
		return err, nil
	}
	
	if todo.CreatedBy != useruid{
		return errors.New("Attempted to get the task belonging to another user"), nil
	}
	return nil, todo
}
