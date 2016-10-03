package gobump

import "errors"

var PLUGIN_CONFIG_FILENAME = "config.json"
var GOBUMP_MAJOR_VERSION = 0
var GOBUMP_MINOR_VERSION = 1
var GOBUMP_VERSION_STRING = "0.1"


var GOBUMP_API_VERSION = 1

var halt = errors.New("Execution interrupt")
var ERR_NOEXPORT = errors.New("Module export does not exist. ")
var ERR_EXPORTEXISTS = errors.New("Module export exists.")
var ERR_MODNOTFOUND = errors.New("Module not found.")
var ERR_MODALREADYLOADED = errors.New("Module already Loaded.")
