package huasheng

import (
	"net/http"
)

//var templates = template.Must(template.ParseFiles("templates/goodslist.html", "templates/index.html", "templates/login.html",
	//"templates/register.html", "templates/goods.html"))

/*func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}*/

func StartServer() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	for _, handler := range handlers {
		http.HandleFunc(handler.URL, handler.HandlerFunc)
	}
	
	http.ListenAndServe(":8080", nil)
}
