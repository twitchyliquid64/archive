package main


import (
	"github.com/hoisie/web"
	"github.com/cznic/ql"
	"orgops/router"
	"orgops/templates"
	"orgops/data"
	"os/signal"
	"fmt"
	"os"
)

func main() {
	templates.LoadTemplates()
	
	//open database
	db, err := ql.OpenFile("orgdata", &ql.Options{CanCreate: true})
	if err != nil{
		fmt.Println("DB Err on initialisation: ", err)
		return
	}
	data.DB = db
	
	
	
	//watch for kill signals so we have time to close down the database
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		
		ec := db.Close()
		switch {
		case ec != nil && err != nil:
			fmt.Println(ec)
		case ec != nil:
			err = ec
		}
		fmt.Println("Recieved", s, "signal. Database closed with error:", err)
		os.Exit(0)
	}()
	
	//make sure the database has all the tables it needs
	err = data.SetupDatabase()
	if err != nil{
		fmt.Println("DB Err on validation: ", err)
		return
	}
	
	//register web interface handlers
	web.Get("/install", router.Install)
	web.Post("/install", router.InstallHandler)
	web.Get("/", router.Index)
	web.Get("/login", router.Login)
	web.Get("/user/logout/", router.Logout)
	web.Post("/login", router.LoginHandler)
	web.Get("/user/setting", router.Settings)
	web.Post("/user/setting/password", router.ChangePasswordHandler)
	web.Get("/user/setting/2factor", router.Settings_2f)
	web.Get("/user/setting/account", router.Settings_acctdetails)
	web.Get("/user/delete", router.Settings_delete)
	web.Get("/user/tasks/create", router.Todo_newtodo)
	web.Post("/user/tasks/create", router.Todo_newtodoHandler)
	web.Get("/user/tasks/complete/(.*)", router.Todo_setcompleteHandler)
	web.Get("/user/tasks/(.*)", router.Todo_viewtodo)
	web.Get("/user/events/create", router.Events_newevent)
	web.Post("/user/events/create", router.Events_neweventHandler)
	web.Get("/user/events/(.*)", router.Events_viewevent)
	web.Get("/user/files/(.*)", router.Files_show)
	web.Post("/user/uploader/(.*)", router.Files_uploader)
	web.Get("/user/download/(.*)", router.Files_Get)
	web.Get("/user/delete/(.*)", router.Files_Remove)
	web.Get("/user/rename/(.*)", router.Files_Rename)
	web.Post("/user/rename/(.*)",router.Files_RenameHandler)
	web.Post("/user/events/delete/(.*)", router.Events_deleteevent)
	web.Post("/user/newfolder/(.*)", router.Files_newFolderHandler)
	web.Get("/admin/users/createnewuser", router.Admin_newuser)
	web.Get("/admin/users/(.*)/changepw", router.Admin_changepw)
	web.Post("/admin/users/(.*)/changepw", router.Admin_changepwHandler)
	web.Get("/admin/users/(.*)", router.Admin_viewUser)
	web.Post("/admin/users/(.*)", router.Admin_updateUserHandler)
	web.Get("/admin/users", router.Admin_users)
	web.Post("/admin/docreatenewuser", router.Admin_newuserhandler)
	web.Get("/admin/db/do", router.Admin_datastorage)
	web.Post("/admin/db/do", router.Admin_datastorage)
	web.Get("/admin/repos", router.Admin_datastorage)
	web.Post("/admin/db/sharedfolders/do", router.Admin_newsharedfolderhandler)
	web.Get("/admin/db/shared/delete/(.*)", router.Admin_deletesharedfolderhandler)
	web.Get("/admin", router.Admin)
	web.Get("/API/user/details/(.*)/(.*)", router.API_userdetails)
	web.Get("/API/user/upcoming/(.*)/(.*)", router.API_userupcoming)
	web.Run("localhost:8008")
}
