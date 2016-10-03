package gobump

import (
    "github.com/robertkrimen/otto"
    _ "github.com/robertkrimen/otto/underscore"
)

type BumpCtx struct {
	modules map[string]*Module 
	globalExports map[string]string //function name to package name
	appExports map[string]func(call otto.FunctionCall) otto.Value
}


//Creates a new bump plugin context.
func New()*BumpCtx{
	tmp := &BumpCtx{
		modules: make(map[string]*Module ),
		globalExports: make(map[string]string ),
		appExports: make(map[string]func(call otto.FunctionCall) otto.Value),
	}
	return tmp
}


//Loads a given plugin into the bump context. Errors if the plugin is already loaded or its unique attribute collides with another loaded plugin.
func (ctx *BumpCtx)Load(filepath string)(*Module, error){
	
	_, ok := ctx.modules[filepath]
	if ok{
		return nil, ERR_MODALREADYLOADED
	}
	
	mod, err := moduleFromFile(filepath, ctx)
	if err!=nil{
		return nil, err
	}
	
	ctx.modules[filepath] = mod
	err = ctx.addModuleExportsToContext(mod)
	if err!=nil{
		ctx.Unload(filepath)
		return nil, err
	}
	
	go mod.startup()
	return mod, nil
}


//Unloads a Module, given the filepath used to load it.
func (ctx *BumpCtx)Unload(filepath string)error{
	mod, ok := ctx.modules[filepath]
	if !ok{
		return ERR_MODNOTFOUND
	}
	
	//first lets remove it from the list so it is no longer accessible.
	delete(ctx.modules, filepath)
	
	for _, val := range ctx.globalExports {
		if val == filepath{
			delete(ctx.globalExports, filepath)
		}
	}
	
	//now lets tell it to shutdown
	mod.forceShutdown()
	
	return nil
}

//unloads all Modules from the system.
func (ctx *BumpCtx)Shutdown(){
	for mod, _ := range ctx.modules{
		ctx.Unload(mod)
	}
}


func (ctx *BumpCtx)uniqueExists(name string)bool{
	for _, obj := range ctx.modules{
		if obj.config.Unique == name{
			return true
		}
	}
	return false
}


//returns a map of information about all loaded Modules, where the key is the name of the Module, and the value is a list of attributes about it.
//[0] = Name of the Module as written in the config, garranteeed to be type string
//[1] = Version of the Module as written in the config, garranteeed to be type int
//[2] = Class of the Module as written in the config, garranteeed to be type string
func (ctx *BumpCtx)Modules()map[string][]interface{}{
	var outobj map[string][]interface{} = make(map[string][]interface{})
		
	for name, obj := range ctx.modules{
		if obj.running{
			outobj[name] = []interface{}{obj.config.Name, obj.config.Version, obj.config.Class}
		}
	}
	return outobj
}

//returns a loaded module given its name.
func (ctx *BumpCtx)GetLoaded(name string)(*Module,error){
	mod, ok := ctx.modules[name]
	if !ok{
		return nil, ERR_MODNOTFOUND
	}
	return mod, nil
}

//returns a list of all the functions which are exported for you to .Call().
func (ctx *BumpCtx)Exports()[]string{
	var ret []string
	for fname, _ := range ctx.globalExports{
		ret = append(ret, fname)
	}
	return ret
}


func (ctx *BumpCtx)addModuleExportsToContext(mod *Module)error{
	for _, fname := range mod.config.GlobalExports {
		_, ok := ctx.globalExports[fname]
		if ok {//export exists - this is not okay
			return ERR_EXPORTEXISTS
		}
		ctx.globalExports[fname] = mod.name
	}
	return nil
}

//Calls a globally exported function that a plugin registered with arguments, and returns a value and error.
//For instance, if the module 'et' globally exported a function 'phoneHome' in its config file, assuming it takes
//no arguments you would call ctx.Call('phoneHome').
func (ctx *BumpCtx)Call(fname string, args ...interface{})(otto.Value, error){
	name, ok := ctx.globalExports[fname]
	if !ok {//export does not exist
		return otto.UndefinedValue(), ERR_NOEXPORT
	}
	
	mod, err := ctx.GetLoaded(name)
	if err != nil{
		return otto.UndefinedValue(), err
	}
	
	return mod.externalCall(fname, args...)
}


//Lets you register a go function to be called from JS code. Pass a function name and a function identifier of the form:
//func(call otto.FunctionCall) otto.Value.
//Then, you can call the function from the app namespace. EG: `app.sayHello("Hi");` all in javascript.
func (ctx *BumpCtx)App(fname string, fn func(call otto.FunctionCall) otto.Value){
	if fn != nil{
		ctx.appExports[fname] = fn
	}
}
