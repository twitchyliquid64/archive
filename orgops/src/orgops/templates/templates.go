package templates

import (
	"html/template"
	"path/filepath"
	"strings"
	"sync"
	"fmt"
	"os"
)

var templates *template.Template = template.New("empty")
var updateLock sync.Mutex

//called once on startup
func LoadTemplates(){
	err := filepath.Walk("templates", loadTemplate)
	if err != nil{
		panic(err)
	}
}

//iterator function - called for each file in the templates folder.
func loadTemplate(path string, f os.FileInfo, err error) error {
	updateLock.Lock()
	defer updateLock.Unlock()
	
	if f != nil{
		if strings.HasSuffix(f.Name(), ".tmpl") {
			fmt.Println("Loading template:", path)
			templates, err = templates.ParseFiles(path)
			if err != nil{
				fmt.Println("\t", err)
			}
		}
	}
	return nil
}

//used for debug
func PrintTemplates() {
	for _, t := range templates.Templates() {
		fmt.Println("\t", t.Name())
	}
}

//called to fetch a template.
func GetTemplate(name string)*template.Template {
	updateLock.Lock()
	defer updateLock.Unlock()
	return templates.Lookup(name)
}
