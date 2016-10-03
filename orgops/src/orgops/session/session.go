package session

import (
    "crypto/rand"
)

const (
	KEY_LENGTH = 16
)

var sessionTable *lossytable

func GetSession(key string)map[string]string {
	if key == ""{return nil}
	
	s := sessionTable.getSession(key)
	if s != nil{ return s.attribs }
	return nil
}

func GetUid(key string)int64{
	if key == ""{return 0}
	
	s := sessionTable.getSession(key)
	if s != nil{ return s.UID }
	return 0
}

func CreateSession(useruid int64, attribs map[string]string)string{
	key := randSessionKey(KEY_LENGTH)
	
	sessionTable.insert(useruid, key)
	for k, v := range attribs {
		sessionTable.updateAttrs(key, k, v)
	}
	return key
}


func DeleteSession(key string){
	sessionTable.delete(key)
}


func init(){
	sessionTable = newSessionTable()
}


func randSessionKey(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}

func GetMaxStorage()int{
	return MAX_STORED
}

func GetStored()int{
	return len(sessionTable.Lookup)
}
