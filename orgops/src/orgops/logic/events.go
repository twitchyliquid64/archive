package logic

import(
	"orgops/data"
	"time"
	"strings"
	"strconv"
	"errors"
)

func NewEvent(useruid int64, title, desc, list, date string)error{
	
	//get the components of the input date
	datetime := strings.Split(date, "T")
	spl := strings.Split(datetime[0], "-")
	
	for i:=0;i<len(spl);i++{
		spl[i] = strings.TrimLeft(spl[i], "0")
		if spl[i] == ""{
			spl[i] = "0"
		}
	}
	
	spl2 := strings.Split(datetime[1], ":")
	
	for i:=0;i<len(spl2);i++{
		spl2[i] = strings.TrimLeft(spl2[i], "0")
		if spl2[i] == ""{
			spl2[i] = "0"
		}
	}
	
	year, _ := strconv.ParseInt(spl[0], 0, 32)
	month, _ := strconv.ParseInt(spl[1], 0, 32)
	day, _ := strconv.ParseInt(spl[2], 0, 32)
	hour, _ := strconv.ParseInt(spl2[0], 0, 32)
	min, _ := strconv.ParseInt(spl2[1], 0, 32)
	
	tm := time.Date(int(year), time.Month(int(month)), int(day), int(hour), int(min), 0, 0, time.Local )
	return data.AddEvent(useruid, title, desc, list, tm)
}

func SeeUsrEvent(useruid int64, eventid string)(error, *data.Event){
	eid, err := strconv.ParseInt(eventid, 10, 64)
	if err!=nil{
		return err, nil
	}
	
	err, event := data.GetEvent(eid)
	if err!=nil{
		return err, nil
	}
	
	if event.Owner != useruid{
		return errors.New("Attempted to get the event belonging to another user"), nil
	}
	return nil, event
}


func DeleteEvent(useruid int64, eventid string)error{
	err, event := SeeUsrEvent(useruid, eventid)
	if err!=nil{
		return err
	}
	
	return event.Delete()
}
