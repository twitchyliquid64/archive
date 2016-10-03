package gobump

import (
    "github.com/robertkrimen/otto"
    _ "github.com/robertkrimen/otto/underscore"
	"archive/zip"
	"errors"
	"sync"
	"time"
)


type Module struct{
	zfile *zip.ReadCloser
	name string
	config *ModuleConfig 
	parent *BumpCtx
	sandbox *otto.Otto
	running bool				//set to true as soon as a goroutine is spawned, set to false as soon as it exits
	hasRunStartup bool			//set to true once all the startup scripts have been run
	lock sync.Mutex
}

type ModuleConfig struct {
	Name string							//canonical name to help reference the Module
	Version int							//canonical version to keep track of versioning
	Unique string						//A unique string which can only be had by a single loaded Module. Any other Module being loaded with this unique string will error.
	Class string						//a machine readable string of the catagory which this Module fits into. Intended to be used to group like Modules under a common banner.
	Startup []string					//a list of JS script files in the package which will be run on startup
	Exports map[string][]string			//map[function_name][]parameter_types, where parameter_types = {string|bool|number, any other name simply means custom object}
	GlobalExports []string				//basically where you want your exports to carry forward and be used globally
}


//Calle
func moduleFromFile(filepath string, parent *BumpCtx)(mod *Module, err error){
	mod = &Module{
		name: filepath,
		parent: parent,
		running: false,
		hasRunStartup: false,
	}
	mod.zfile, err = zip.OpenReader(filepath)
	if err != nil{
		return
	}
	
	mod.config, err = readConfig(mod.zfile)//zip.go
	if err != nil{
		return
	}
	
	err = mod.configValidate()
	if err != nil{
		return
	}
	
	mod.sandbox = otto.New()
	mod.sandbox.Interrupt = make(chan func(), 1)
	mod.parent.setupSandboxEnvs(mod.sandbox, mod)
	return
}


func (mod *Module)forceShutdown(){
	if mod.running{
		mod.sandbox.Interrupt <- func() {
            panic(halt)
		}
	}
	mod.executing()//take ownership of the Modules execution context
	mod.zfile.Close()
	mod.parent = nil//speed up garbage collection
	mod.config = nil//speed up garbage collection
}


//now that the Module has a config, it needs to be validated.
func (mod *Module)configValidate()error{
	if mod.config.Unique != ""{
		if mod.parent.uniqueExists(mod.config.Unique){//at this stage the Module hasnt been added to the parent, so it is safe to check that none exist
			return errors.New("There is a Module already loaded with that unique identifier!")
		}
	}
	return nil
}

//hangs until the module has finished running through its startup scripts. Typically, you should call this before you use
//anything the module sets up (such as functions).
func (mod *Module)WaitTillStartupDone(){
	for mod.hasRunStartup == false{//wait for startup to stop running
		time.Sleep(time.Millisecond*10)
	}
}
