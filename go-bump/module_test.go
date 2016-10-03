package gobump

import "testing"


func TestNotExists( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/doesnotexist.zip")
	if err == nil{
		t.Error("Load() did not error, should have fileNotFound")
	}
}

func TestNoConfig( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/noConfig.zip")
	if err == nil{
		t.Error("Load() did not error, should have configNotFound")
	}
}

func TestConfigMalformed( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/configMalformed.zip")
	if err == nil{
		t.Error("Load() did not error, should have configMalformed")
	}
}


func TestEmpty( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/blank.zip")
	if err != nil{
		t.Error("Load() errorred on empty module:", err.Error())
	}
}
