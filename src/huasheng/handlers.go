package huasheng

import (
	"net/http"
	"regexp"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "index.html", map[string]interface{}{})
}

/*func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "login", nil)
	} else {
		user := User{Name: r.FormValue("username"), Password: r.FormValue("password")}
		b, err := user.Login()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if b {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Error(w, "用户名或者密码错误。。。", http.StatusInternalServerError)
		}

	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "register", nil)
	} else {
		user := User{Name: r.FormValue("username"), Email: r.FormValue("email"), Password: r.FormValue("password")}
		b, err := user.IsExist()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if b {
			http.Error(w, "该用户名已注册。。。", http.StatusInternalServerError)
		} else {
			fmt.Println(user.Name)
			err := user.Create()
			if err != nil {
				fmt.Println(user.Name)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Println(user.Name)
			http.Redirect(w, r, "/", http.StatusFound)
		}

	}
}*/

func goodslistHandler(w http.ResponseWriter, r *http.Request) {
	goodes, err := loadGoodses()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	renderTemplate(w, r, "goodslist.html", map[string]interface{}{"goods": goodes})
}

var validPath = regexp.MustCompile("^/(goods)/(中原花生油大罐装|中原花生油小罐装)$")

func regexpGoodsHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
func goodsHandler(w http.ResponseWriter, r *http.Request, name string) {
	g, err := loadGoods(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	renderTemplate(w, r, "goods.html", map[string]interface{}{"goods": g})
}
