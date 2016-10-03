package files

import (
	"path/filepath"
	"io/ioutil"
	"strings"
	"fmt"
	"os"
	"io"
)

func Remove(username, path string)error{
	return os.RemoveAll(filepath.Join(getUserHome(username), sanitizeInputPath(path)))
}

func NewFolder(username, path, lower, newname string)error{
	return os.Mkdir(filepath.Join(getUserHome(username), sanitizeInputPath(path), sanitizeInputPath(newname)), 0777)
}



func Rename(username, path, lower, newname string)error{
	oldPath := filepath.Join(getUserHome(username), sanitizeInputPath(path))
	newPath := filepath.Join(getUserHome(username), sanitizeInputPath(lower), sanitizeInputPath(newname))
	
	return os.Rename(oldPath, newPath)
}

func IsFile(username, path string)(bool, error){
	fullPath := filepath.Join(getUserHome(username), sanitizeInputPath(path))
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


func Read(username, path string)([]byte, error){
	fullPath := filepath.Join(getUserHome(username), sanitizeInputPath(path))
	return ioutil.ReadFile(fullPath)
}


func DoUpload(username, path, filename string, file io.ReadCloser)error{
	fullPath := filepath.Join(getUserHome(username), sanitizeInputPath(path), sanitizeInputPath(filename))
	fmt.Println("Doing upload to:" + fullPath)
	
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


//shared folder name is passed as parameter 'username' when isSharedFolder is true.
func GetFileList(username, path string, isSharedFolder bool)([]map[string]interface{}, error){
	fullPath := filepath.Join(getUserHome(username), sanitizeInputPath(path))
	
	if isSharedFolder{
		fullPath = filepath.Join(getSharedFolderHome(username), sanitizeInputPath(path))
	}
	
	fileData, err := ioutil.ReadDir(fullPath)
	if err != nil{
		return nil, err
	}
	
	var ret []map[string]interface{}
	
	for _, file := range fileData{
		ret = append(ret, map[string]interface{}{
			"Name": file.Name(),
			"Size": sizeString(file.Size()),
			"IsDir": file.IsDir(),
			"Modified": file.ModTime().Format("Jan 2, 2006 3:04pm"),
		})
	}
	return ret, nil
}

func getUserHome(username string)string{
	return "userdata/"+username
}


func sanitizeInputPath(path string)string{
	return strings.Replace(filepath.Clean(path), "..", "", -1)
}

func boolToString(in bool)string{
	if in{
		return "true"
	}
	return "false"
}

func sizeString(size int64)string{
	var suffixes = []string{"B","KB","MB","GB","TB","PB","EB", "ZB", "YB"}
	
	for _, suffix := range suffixes{
		if size < 1024{
			return fmt.Sprintf("%d %s", size, suffix)
		}
		size /= 1024
	}
	return "ERR"
}
