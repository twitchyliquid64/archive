package files

import (
	"path/filepath"
	"orgops/data"
	"io/ioutil"
	"strings"
	"errors"
	"os"
	"io"
)


func NewSharedFolder(name string)error{
	return os.Mkdir(getSharedFolderHome(name), 0777)
}



func getSharedFolderHome(name string)string{
	return "shareddata/"+name
}


func SharedFolderIsFile(s map[string]string, path string)(bool, error){
	spl := strings.Split(path, "/")
	err, isSharedFolder := data.CanAccessSharedFolder(spl[0], s["username"])
	if err != nil {
		return false, err
	}
	
	if isSharedFolder{
		fullPath := filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(strings.Join(spl[1:], "/")))
		file, err := os.OpenFile(fullPath, os.O_WRONLY, 0777)
		if err != nil {
			return false, err
		}
		defer file.Close()
		
		fi, err := file.Stat()
		if err != nil {
			return false, err
		}
		
		return !fi.IsDir(), nil
	}
	return false, nil
}

func SharedFolderRead(s map[string]string, path string)([]byte, error){
	spl := strings.Split(path, "/")
	err, isSharedFolder := data.CanAccessSharedFolder(spl[0], s["username"])
	if err != nil {
		return nil, err
	}
	if !isSharedFolder{
		return nil, errors.New("Permission denied")
	}
	fullPath := filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(strings.Join(spl[1:], "/")))
	return ioutil.ReadFile(fullPath)
}


//assumes that access control check as already been done by calling function.
func SharedFolderDoUpload(username, path, filename string, file io.ReadCloser)error{
	spl := strings.Split(path, "/")
	fullPath := filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(strings.Join(spl[1:], "/")), sanitizeInputPath(filename))
	
	defer file.Close()
	dst, err := os.OpenFile(fullPath, os.O_CREATE | os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}

//assumes that access control check as already been done by calling function.
func SharedFolderRemove(username, path string)error{
	spl := strings.Split(path, "/")
	fullPath := filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(strings.Join(spl[1:], "/")))
	
	return os.RemoveAll(fullPath)
}

//assumes that access control check as already been done by calling function.
func SharedFolderRename(username, path, lower, newname string)error{
	spl := strings.Split(path, "/")
	oldPath := filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(strings.Join(spl[1:], "/")))
	newPath := filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(getLower(strings.Join(spl[1:], "/"))), sanitizeInputPath(newname))
	
	return os.Rename(oldPath, newPath)
}

//assumes that access control check as already been done by calling function.
func SharedFolderNewFolder(username, path, lower, newname string)error{
	spl := strings.Split(path, "/")
	return os.Mkdir(filepath.Join(getSharedFolderHome(spl[0]), sanitizeInputPath(strings.Join(spl[1:], "/")), sanitizeInputPath(newname)), 0777)
}


func getLower(path string)string{
	spl := strings.Split(strings.TrimRight(path, "/"), "/")
	lower := ""
	if len(spl) < 2{
		return ""
	}
	
	for _, level := range spl[:len(spl)-1]{
		lower += level+"/"
	}
	return lower
}
