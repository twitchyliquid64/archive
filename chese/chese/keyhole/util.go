package keyhole

import (
	"html/template"
	"net/http"
	"path"
)

// renderPage renders the template at the specific path.
func renderPage(pagePath string, pageData interface{}, rw http.ResponseWriter) error {
	t, err := template.New("").Delims("{!{", "}!}").ParseFiles(pagePath)
	if err != nil {
		http.Error(rw, "Internal Server Error", 500)
		return err
	}
	return t.ExecuteTemplate(rw, path.Base(pagePath), pageData)
}
