package gobump

import (
	"fmt"
	"time"
    "github.com/robertkrimen/otto"
    _ "github.com/robertkrimen/otto/underscore"
)


//This function SHOULD be run in a goroutine
func (mod *Module)startup(){
	
	if mod.hasRunStartup{//can only run once
		return
	}
	
	mod.executing()
	defer mod.leaving()
	
	for _, scriptname := range mod.config.Startup{
		data, err := readFile(scriptname, mod.zfile)
		
		if err != nil{
			fmt.Println(err)
			continue
		}
		
		val, err := mod.sandbox.Run(data)
		if err != nil {
			fmt.Println(err, val)
		}
	}
	
	mod.hasRunStartup = true
}

func (mod *Module)executing(){
	mod.lock.Lock()
	mod.running = true
}

func (mod *Module)leaving(){
	mod.lock.Unlock()
	mod.running = false
	if caught := recover(); caught != nil {
		if caught == halt {
			fmt.Println("("+mod.name+") Execution interrupted")
		}
	}
}


func (ctx *BumpCtx)setupSandboxEnvs(sandbox *otto.Otto, mod *Module){
	
	//set go-bump version
	mod.sandbox.Set("GOBUMP_MINOR_VERSION", GOBUMP_MINOR_VERSION)
	mod.sandbox.Set("GOBUMP_MAJOR_VERSION", GOBUMP_MAJOR_VERSION)
	mod.sandbox.Set("GOBUMP_VERSION_STRING", GOBUMP_VERSION_STRING)
	
	//setup the general API
	api, _ := mod.sandbox.Object(`api = {}`)
	api.Set("version", GOBUMP_API_VERSION)
	api.Set("modulesLoaded", func(call otto.FunctionCall) otto.Value {
		result, _ := mod.sandbox.ToValue(ctx.Modules())
		return result
	})
	
	api.Set("moduleExports", func(call otto.FunctionCall) otto.Value {
		modS, _ := call.Argument(0).ToString()
		modO, ok := ctx.modules[modS]
		if !ok{
			return otto.UndefinedValue()
		}

		result, _ := mod.sandbox.ToValue(modO.config.Exports)
		return result
	})
	
	api.Set("moduleCall", func(call otto.FunctionCall) otto.Value {
		modS, _ := call.Argument(0).ToString()
		modO, ok := ctx.modules[modS]
		if !ok{
			return otto.UndefinedValue()
		}
		mod.leaving()//we are calling remote so we can free up this Module for execution.
		defer mod.executing()//on return, we are returning to this Module so we should grab the lock on execution.
		return modO.remoteModuleCall(call)
	})
	
	api.Set("moduleByClass", func(call otto.FunctionCall) otto.Value {
		var resultset []map[string]interface{}
		class, _ := call.Argument(0).ToString()
		
		for _, val := range ctx.modules{
			if class == val.config.Class{
				resultset = append(resultset, map[string]interface{}{
					"name": val.config.Name,
					"version": val.config.Version,
					"class": val.config.Class,
					"exports": val.config.Exports,
				})
			}
		}

		result, _ := mod.sandbox.ToValue(resultset)
		return result
	})
	
	api.Set("sleep", func(call otto.FunctionCall) otto.Value {
		duration, _ := call.Argument(0).ToInteger()
		time.Sleep(time.Millisecond * time.Duration(duration))
		return otto.UndefinedValue()
	})
	
	
	//setup application exports
	app, _ := mod.sandbox.Object(`app = {}`)
	for fname, fn := range ctx.appExports {
		app.Set(fname, fn)
	}
	
	
	mod.setupSandboxEnvs()
}


func (mod *Module)setupSandboxEnvs(){

}

//called when another Module is calling a JS function in this Module.
func (mod *Module)remoteModuleCall(call otto.FunctionCall) otto.Value {
	mod.executing()
	defer mod.leaving()
	var args []interface{}
	pos := 2
	for {
		args = append(args, call.Argument(pos))
		if call.Argument(pos).IsUndefined(){
			break
		}
		pos++
	}
	
	fname, _ := call.Argument(1).ToString()
	val, err := mod.sandbox.Call(fname, nil, args...)
	if err!=nil{
		fmt.Println(err)
	}
	return val
}



//called from context to execute a function in the module context and return a result.
func (mod *Module)externalCall(fname string, args ...interface{})(otto.Value,error){
	mod.executing()
	defer mod.leaving()
	return mod.sandbox.Call(fname, nil, args...)
}
