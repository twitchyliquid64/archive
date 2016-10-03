package gobump

import "testing"
import "strconv"

func TestApi( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/testApi.zip")
	if err != nil{
		t.Error("Load() errorred: ", err)
	}
	
	ctx.modules["tests/testApi.zip"].WaitTillStartupDone()
	
	//api.version
	value, err := ctx.modules["tests/testApi.zip"].sandbox.Get("TESTING_HARNESS")
	if err != nil{
		t.Error("sandbox.Get() errorred: ", err)
	}

	val, err := value.ToString()
	if err != nil{
		t.Error("value.ToString() errorred: ", err)
	}
	if val != (GOBUMP_VERSION_STRING+strconv.Itoa(GOBUMP_API_VERSION)){
		t.Error("sandbox.Get() unexpected: ", val, (GOBUMP_VERSION_STRING+strconv.Itoa(GOBUMP_API_VERSION)))
	}
	
	
	//api.ModulesLoaded
	value, err = ctx.modules["tests/testApi.zip"].sandbox.Get("TESTING_HARNESS_MODULES")
	if err != nil{
		t.Error("sandbox.Get() errorred: ", err)
	}

	val, err = value.ToString()
	if err != nil{
		t.Error("value.ToString() errorred: ", err)
	}
	if val != "{\"tests/testApi.zip\":[\"" + ctx.modules["tests/testApi.zip"].config.Name + "\"," + strconv.Itoa(ctx.modules["tests/testApi.zip"].config.Version) + ",\"" + ctx.modules["tests/testApi.zip"].config.Class + "\"]}"{
		t.Error("sandbox.Get() unexpected: ", val)
	}
	
	
	
	//api.ModuleExports
	value, err = ctx.modules["tests/testApi.zip"].sandbox.Get("TESTING_HARNESS_MODLIST")
	if err != nil{
		t.Error("sandbox.Get() errorred: ", err)
	}

	val, err = value.ToString()
	if err != nil{
		t.Error("value.ToString() errorred: ", err)
	}
	if val != "{\"testExport\":[\"int\",\"string\",\"bool\"]}"{
		t.Error("sandbox.Get() unexpected: ", val)
	}
	
	ret := ctx.Exports()
	if len(ret) != 1 || ret[0] != "testExport"{
		t.Error("ctx.Exports() unexpected: ", ret, len(ret))
	}
}


func TestApiModuleCallSelf( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/testSelfCall.zip")
	if err != nil{
		t.Error("Load() errorred: ", err)
	}
	
	ctx.modules["tests/testSelfCall.zip"].WaitTillStartupDone()
	
	//api.ModuleExports
	value, err := ctx.modules["tests/testSelfCall.zip"].sandbox.Get("TESTING_HARNESS_RETVAL")
	if err != nil{
		t.Error("sandbox.Get() errorred: ", err)
	}

	val, err := value.ToInteger()
	if err != nil{
		t.Error("value.ToInteger() errorred: ", err)
	}
	if val != 888{
		t.Error("value.ToInteger() returned unexpected: ", val, " expecting 888")
	}
}



func TestApiModuleByClass( t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/testModuleByClass.zip")
	if err != nil{
		t.Error("Load() errorred: ", err)
	}
	
	ctx.modules["tests/testModuleByClass.zip"].WaitTillStartupDone()
	
	//api.ModuleExports
	value, err := ctx.modules["tests/testModuleByClass.zip"].sandbox.Get("TESTING_HARNESS_BYCLASS")
	if err != nil{
		t.Error("sandbox.Get() errorred: ", err)
	}

	val, err := value.ToString()
	if err != nil{
		t.Error("value.ToString() errorred: ", err)
	}
	if val != "[{\"class\":\"testClassName\",\"exports\":{},\"name\":\"testClass\",\"version\":1}]"{
		t.Error("value.ToString() returned unexpected: ", val, " expecting [{\"class\":\"testClassName\",\"exports\":{},\"name\":\"testClass\",\"version\":1}]")
	}
}

func TestExportCollision(t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/testCollision.zip")
	if err == nil{
		t.Error("Load() should have errorred with a collision.", err)
	}
}


func TestExport(t *testing.T){
	ctx := New()
	_, err := ctx.Load("tests/testExport.zip")
	if err != nil{
		t.Error("Load() errorred: ", err)
	}
	
	ctx.modules["tests/testExport.zip"].WaitTillStartupDone()
	value, err := ctx.Call("testExport", "teststr", 3)
	if err != nil{
		t.Error("Call() should not have errorred:", err)
	}
	val, err := value.ToString()
	if err != nil{
		t.Error("value.ToString() errorred: ", err)
	}
	if val != "teststr 9"{
		t.Error("Call() returned unexpected string:", val, "expected teststr 9")
	}
}
