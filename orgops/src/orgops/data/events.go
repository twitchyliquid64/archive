package data

import (
	"github.com/cznic/ql"
	"time"
)

type Event struct {
	UID int64
	Owner int64
	Created time.Time
	
	Title string
	Description string
	Time time.Time

	List string
}


var addEventQry = ql.MustCompile("BEGIN TRANSACTION; INSERT INTO event(uid, title, description, list, created, eventtime) VALUES ($1,$2,$3,$4,now(),$5); COMMIT;")

func AddEvent(useruid int64, title, description, list string, eventTime time.Time)error{
	_, _, err := DB.Execute(ql.NewRWCtx(), addEventQry, useruid, title, description, list, eventTime)
	return err
}

var getUsrEventQry = ql.MustCompile("SELECT id(), title, description, eventtime, list, created FROM event WHERE uid == $1 && eventtime > now() ORDER BY eventtime ASC LIMIT 20;")

func EventByUserUid(uid int64)(error, []Event){
	var ret []Event
	results, _, err := DB.Execute(ql.NewRWCtx(), getUsrEventQry, uid)
	if err != nil{
		return err, nil
	}
	
	eventData, err := results[0].Rows(2048,0)//max 2048 todo's
	if err != nil{
		return err, nil
	}
	
	for _, event := range eventData {
		ret = append(ret, Event{
			UID: event[0].(int64),
			Title: event[1].(string),
			Description: event[2].(string),
			Time: event[3].(time.Time),
			List: event[4].(string),
			Created: event[5].(time.Time),
			Owner: uid,//users uid, not task uid
		})
	}
	return nil, ret
}




var getEventQry = ql.MustCompile("SELECT title, description, created, list, eventtime, uid FROM event WHERE id() == $1;")
func GetEvent(uid int64)(error, *Event){
	results, _, err := DB.Execute(ql.NewRWCtx(), getEventQry, uid)
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
	
	return nil, &Event{
			UID: uid,
			Title: todo[0].(string),
			Description: todo[1].(string),
			Time: todo[4].(time.Time),
			Created: todo[2].(time.Time),
			List: todo[3].(string),
			Owner: todo[5].(int64),
	}
}


var delEventQry = ql.MustCompile("BEGIN TRANSACTION; DELETE FROM event WHERE id() == $1; COMMIT;")
func (inst *Event)Delete()error{
	_, _, err := DB.Execute(ql.NewRWCtx(), delEventQry, inst.UID)
	return err
}
