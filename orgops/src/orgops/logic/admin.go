package logic

import (
	"orgops/session"
	"orgops/files"
	"orgops/data"
	"time"
	"runtime"
	"fmt"
)

var startTime = time.Now()

var sysStatus struct {
	Uptime       string
	NumGoroutine int
	// Sessuib store statistics
	SessionsStored	int
	SessionSlots	int

	// General statistics.
	MemAllocated int64 // bytes allocated and still in use
	MemTotal     int64 // bytes allocated (even if freed)
	MemSys       int64 // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    int64 // bytes allocated and still in use
	HeapSys      int64 // bytes obtained from system
	HeapIdle     int64 // bytes in idle spans
	HeapInuse    int64 // bytes in non-idle span
	HeapReleased int64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  int64 // bootstrap stacks
	StackSys    int64
	MSpanInuse  int64 // mspan structures
	MSpanSys    int64
	MCacheInuse int64 // mcache structures
	MCacheSys   int64
	BuckHashSys int64 // profiling bucket hash table
	GCSys       int64 // GC metadata
	OtherSys    int64 // other system allocations

	// Garbage collector statistics.
	NextGC       int64 // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}

func Status()interface{} {
	sysStatus.Uptime = startTime.Format("Jan 2, 2006 at 3:04pm (MST)")

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	sysStatus.NumGoroutine = runtime.NumGoroutine()
	
	sysStatus.SessionsStored = session.GetStored()
	sysStatus.SessionSlots = session.GetMaxStorage()

	sysStatus.MemAllocated = (int64(m.Alloc)) / 1024
	sysStatus.MemTotal = (int64(m.TotalAlloc)) / 1024
	sysStatus.MemSys = (int64(m.Sys)) / 1024
	sysStatus.Lookups = m.Lookups
	sysStatus.MemMallocs = m.Mallocs
	sysStatus.MemFrees = m.Frees

	sysStatus.HeapAlloc = (int64(m.HeapAlloc))
	sysStatus.HeapSys = (int64(m.HeapSys))
	sysStatus.HeapIdle = (int64(m.HeapIdle))
	sysStatus.HeapInuse = (int64(m.HeapInuse))
	sysStatus.HeapReleased = (int64(m.HeapReleased))
	sysStatus.HeapObjects = m.HeapObjects

	sysStatus.StackInuse = (int64(m.StackInuse))
	sysStatus.StackSys = (int64(m.StackSys))
	sysStatus.MSpanInuse = (int64(m.MSpanInuse))
	sysStatus.MSpanSys = (int64(m.MSpanSys))
	sysStatus.MCacheInuse = (int64(m.MCacheInuse))
	sysStatus.MCacheSys = (int64(m.MCacheSys))
	sysStatus.BuckHashSys = (int64(m.BuckHashSys))
	sysStatus.GCSys = (int64(m.GCSys))
	sysStatus.OtherSys = (int64(m.OtherSys))

	sysStatus.NextGC = (int64(m.NextGC))
	sysStatus.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	sysStatus.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	sysStatus.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	sysStatus.NumGC = m.NumGC
	
	return sysStatus
}



func IsAdmin(usr string)bool{
	err, user := data.GetUserByName(usr)
	if err!=nil{
		fmt.Println("DB Error:", err)
		return false
	}
	if user == nil{return false}
	
	return chkStrInStrSlice("admin", user.Groups)
}

func GetUserList()(error,[]*data.User){
	return data.GetUsers()
}

func Admin_UpdateUser(user, addgroup, removegroup string)(err error){
	
	if addgroup != ""{
		err = data.AddUserToGroupByUsername(user, addgroup)
		if err!=nil{ return err }
	}
	if removegroup != ""{
		err = data.DelUserFromGroupByUsername(user, removegroup)
		if err!=nil{ return err }
	}
	return nil
}


func Admin_ChangePw(user, passwd string)error{
	return ChangePassword(user, "", passwd, passwd, true)
}

func chkStrInStrSlice(str string, slice []string)bool{
	for _, row := range slice{
		if row == str{
			return true
		}
	}
	return false
}

func Admin_NewUser(user, pass string)error{
	return data.NewUser(&data.User{
		Username: 	user,
		Passhash: 	shaHash(user, pass),
		Groups:		[]string{},
	})
}

func SetupAdminAccount(user, pass string)bool{
	if data.UsrCount() == 0 {
		err := data.NewUser(&data.User{
			Username: 	user,
			Passhash: 	shaHash(user, pass),
			Groups:		[]string{"admin"},
		})
		return err == nil//true on success
	}
	return false
}


func Admin_DoQuery(input string)(err error, fields []string, rows [][]interface{}){
	err, rec := data.Query(input)
	if err!=nil{
		return
	}
	if rec == nil{
		return
	}
	
	if len(*rec) == 0{
		return
	}
	
	fields, err = (*rec)[0].Fields()
	rows, err = (*rec)[0].Rows(1024, 0)
	return
}


func Admin_NewSharedFolder(foldername, groupname string)error{
	fmt.Println("logic.Admin_NewSharedFolder() :: ", foldername, groupname)
	err := files.NewSharedFolder(foldername)
	if err != nil {
		fmt.Println("\tfiles.NewSharedFolder() returned error ", err)
		return err
	}
	err = data.AddSharedFolder(foldername, groupname)
	return err
}

func Admin_DeleteSharedFolder(groupname string)error{
	return data.DelSharedFolder(groupname)
}
