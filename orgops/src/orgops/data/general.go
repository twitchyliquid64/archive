package data

import (
	"github.com/cznic/ql"
	"crypto/rand"
)

var DB *ql.DB


func Query(input string)(error, *[]ql.Recordset){
	rs, _, err := DB.Run(ql.NewRWCtx(), input)
	if err!=nil{
		return err, nil
	}
	
	return nil, &rs
}



func randString(n int) string {
    const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}
