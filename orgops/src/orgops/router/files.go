package router

import (
	"github.com/hoisie/web"
	"orgops/templates"
	"orgops/session"
	"path/filepath"
	"orgops/files"
	"orgops/data"
	"strings"
	"fmt"
)

func Files_Get(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		data, err := files.Read(s["username"], path)
		if err != nil{
			data, _ = files.SharedFolderRead(s, path)
			fmt.Println(err)
		}
		
		
		ctx.ContentType("application/octet-stream")
		ctx.SetHeader("Content-Disposition", "attachment", true)
		ctx.ResponseWriter.Write(data)
	}
}


func Files_Rename(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	lower := getLower(path)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		renderTemplate(ctx, "file_rename.tmpl", FilePageData(s, nil, nil, path, strings.Split(path, "/"), lower, nil))
	}
}

func Files_newFolderHandler(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	lower := getLower(path)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		spl := strings.Split(path, "/")
		err, isSharedFolder := data.CanAccessSharedFolder(spl[0], s["username"])
		if err != nil {
			fmt.Println("data.CanAccessSharedFolder() returned error:", err)
		}
		if isSharedFolder{
			err = files.SharedFolderNewFolder(s["username"], path, lower, ctx.Params["fname"])
		}else{
			err = files.NewFolder(s["username"], path, lower, ctx.Params["fname"])
		}
		if err != nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/user/files/"+path)
	}
}


func Files_RenameHandler(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	lower := getLower(path)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		isSharedFile, err := files.SharedFolderIsFile(s, path)
		if isSharedFile{
			err = files.SharedFolderRename(s["username"], path, lower, ctx.Params["fname"])
		}else{
			err = files.Rename(s["username"], path, lower, ctx.Params["fname"])
		}
		if err != nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/user/files/"+lower)
	}
}




func Files_Remove(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	lower := getLower(path)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		spl := strings.Split(path, "/")
		err, isSharedFolder := data.CanAccessSharedFolder(spl[0], s["username"])
		if err != nil {
			fmt.Println("data.CanAccessSharedFolder() returned error:", err)
		}
		if isSharedFolder{
			err = files.SharedFolderRemove(s["username"], path)
		}else{
			err = files.Remove(s["username"], path)
		}
		if err != nil{
			fmt.Println(err)
		}
		ctx.Redirect(302, "/user/files/"+lower)
	}
}


func Files_show(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	lower := getLower(path)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		isFile, err := files.IsFile(s["username"], path)
		isSharedFile, err2 := files.SharedFolderIsFile(s, path)
		if (err != nil) || (err2 != nil){
			fmt.Println("files.IsFile() returned", isFile, "error:", err)
			fmt.Println("files.SharedFolderIsFile() returned", isSharedFile, "error:", err2)
		}
		if isFile || isSharedFile{
			ext := strings.ToLower(filepath.Ext(path))
			var data []byte
			var err error
			if isSharedFile{
				data, err = files.SharedFolderRead(s, path)
			}else{
				data, err = files.Read(s["username"], path)
			}
			if err != nil{
				fmt.Println(err)
			}
			
			//if we have a renderer, load the renderer. Else, download the file.
			if(templates.GetTemplate("filerenderer_"+ext+".tmpl") != nil) {
				renderTemplate(ctx, "filerenderer_"+ext+".tmpl", FilePageData(s, string(data), err, path, strings.Split(path, "/"), lower, nil))
			}else{
				ctx.ContentType("application/octet-stream")
				ctx.SetHeader("Content-Disposition", "attachment", true)
				ctx.ResponseWriter.Write(data)
			}
		}else{
			
			fileList, err := files.GetFileList(s["username"], path, false)
			
			spl := strings.Split(path, "/")
			if err != nil{//if there was an error, lets see if its actually a shared folder
				err2, isSharedFolder := data.CanAccessSharedFolder(spl[0], s["username"])
				if err2 != nil {
					fmt.Println("data.CanAccessSharedFolder(",spl[0], s["username"],") ::", err2, isSharedFolder)
				}
				if isSharedFolder{
					fileList, err = files.GetFileList(spl[0], strings.Join(spl[1:], "/"), true)
				}
			}
			
			var usrShared []string
			var err3 error
			if (path == "") || (path == "/"){//if we are in root, show the shared folders
				err3, usrShared = data.GetUserFolders(session.GetUid(getSessionKey(ctx)))
				if err3!=nil{
					fmt.Println("data.GetUserFolders() err ::", err3)
				}
			}
			
			renderTemplate(ctx, "filebrowser.tmpl", FilePageData(s, fileList, err, path, strings.Split(path, "/"), lower, usrShared))
		}
	}
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


func Files_uploader(ctx *web.Context, path string) {
	s := sessionDetails(ctx)
	
	if s == nil{
		ctx.Redirect(302, "/login")
	}else{
		
		ctx.Request.ParseMultipartForm(10 * 1024 * 1024)
		form := ctx.Request.MultipartForm
		for _, fhd := range form.File["files"]{
			filename := fhd.Filename
			
			file, err := fhd.Open()
			if err != nil {
				fmt.Println(err)
				ctx.Abort(500, "Internal server error")
				return
			}
			
			spl := strings.Split(path, "/")
			err, isSharedFolder := data.CanAccessSharedFolder(spl[0], s["username"])
			if err != nil {
				fmt.Println("data.CanAccessSharedFolder() returned error:", err)
			}
			if isSharedFolder{
				err = files.SharedFolderDoUpload(s["username"], path, filename, file)
			}else{
				err = files.DoUpload(s["username"], path, filename, file)
			}
			if err != nil {
				fmt.Println(err)
				ctx.Abort(500, "Internal server error - please notify your administrator. Error details recorded in server log.")
				return
			}
		}
		ctx.Redirect(302, "/user/files/"+path)
	}
}
