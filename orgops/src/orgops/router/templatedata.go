package router



type TemplateData struct {
	IsDashPage bool
	IsHelpPage bool
	IsSettingPage bool
	IsAdminPage bool
	IsFilePage bool
	
	UserIsAdmin bool
	
	Username string
	Error string
	AffirmativeMessage string
	
	Groups []string //only used on the account page - specialcase
	
	PageSpecificData interface{}//data for specifically a kind of page
	PageSpecificData2 interface{}//data for specifically a kind of page
	PageSpecificData3 interface{}//data for specifically a kind of page
	PageSpecificData4 interface{}//data for specifically a kind of page
	PageSpecificData5 interface{}//this is getting kind of rediculous. No more.
}

type UserUpdateData struct { // used on the updateuser page
	Username string
	IsAdmin bool
	Groups []string
}


func FilePageData(sessiondata map[string]string, files interface{}, err error, path string, paths []string, oneDownPath string, sharedFolders []string)TemplateData{
	errstr := ""
	if err != nil{errstr = err.Error()}
	
	return TemplateData{
		IsFilePage: true,
		Username: sessiondata["username"],
		UserIsAdmin: isAdmin(sessiondata),
		PageSpecificData: files,
		Error: errstr,
		PageSpecificData2: path,
		PageSpecificData3: paths,
		PageSpecificData4: oneDownPath,
		PageSpecificData5: sharedFolders,
	}
}


func DashboardPageData(sessiondata map[string]string, todos, events interface{})TemplateData{
	return TemplateData{
		IsDashPage: true,
		Username: sessiondata["username"],
		UserIsAdmin: isAdmin(sessiondata),
		PageSpecificData: todos,
		PageSpecificData2: events,
	}
}


func DataStoragePageData(input string, headings []string, data [][]interface{}, sharedFolders []*[]string)interface{}{

	return struct{
		Query			string
		Headings 		[]string
		Data			[][]interface{}
		SharedFolders	[]*[]string
	}{
		Query:			input,
		Headings:		headings,
		Data:			data,
		SharedFolders: 	sharedFolders,
	}
}

func AdminPageData(sessiondata map[string]string, stats interface{}, err error)TemplateData{
	errstr := ""
	if err != nil{errstr = err.Error()}
	
	return TemplateData{
		IsAdminPage: true,
		PageSpecificData: stats,
		Username: sessiondata["username"],
		UserIsAdmin: isAdmin(sessiondata),
		Error: errstr,
	}
}


func UserSettingsPageData(sessiondata map[string]string, err error, amsg string, grp []string, apikey string)TemplateData{
	errstr := ""
	if err != nil{errstr = err.Error()}
	
	return TemplateData{
		IsSettingPage: true,
		Username: sessiondata["username"],
		Error: errstr,
		AffirmativeMessage: amsg,
		Groups: grp,
		UserIsAdmin: isAdmin(sessiondata),
		PageSpecificData: apikey,
	}
}
