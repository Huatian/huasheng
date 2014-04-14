package huasheng

import (
	"bytes"
	"fmt"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
	"github.com/jimmykuu/wtforms"
	"runtime"
	"strings"
	"crypto/md5"
	"io"
)

var (
	DB *mgo.Database
	store       *sessions.CookieStore
	utils       *Utils
)

var funcMaps = template.FuncMap{
	"gravatar": func(email string, size uint16) string {
		h := md5.New()
		io.WriteString(h, email)
		return fmt.Sprintf("http://www.gravatar.com/avatar/%x?s=%d", h.Sum(nil), size)
	},
}

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println("MongoDB连接失败:", err.Error())
		os.Exit(1)
	}

	session.SetMode(mgo.Monotonic, true)

	DB = session.DB("huasheng")
	
	store = sessions.NewCookieStore([]byte("huatian"))
	
	utils = &Utils{}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, file string, data map[string]interface{}) {
	_, isPresent := data["signout"]

	// 如果isPresent==true，说明在执行登出操作
	if !isPresent {
		// 加入用户信息
		user, ok := currentUser(r)

		if ok {
			data["username"] = user.Username
			data["isSuperUser"] = user.IsSuperuser
			data["email"] = user.Email
		}
	}

	data["utils"] = utils

	//data["analyticsCode"] = analyticsCode
	//data["shareCode"] = shareCode
	//data["staticFileVersion"] = Config.StaticFileVersion
	data["goVersion"] = runtime.Version()

	_, ok := data["active"]
	if !ok {
		data["active"] = ""
	}

	page := parseTemplate(file, data)
	w.Write(page)
}

func parseTemplate(file string, data map[string]interface{}) []byte {
	var buf bytes.Buffer

	t, err := template.ParseFiles("templates/base.html", "templates/"+file)
	if err != nil {
		panic(err)
	}
	t = t.Funcs(funcMaps)
	err = t.Execute(&buf, data)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

type Utils struct {
}

// 没有http://开头的增加http://
func (u *Utils) Url(url string) string {
	if strings.HasPrefix(url, "http://") {
		return url
	}

	return "http://" + url
}

func (u *Utils) Index(index int) int {
	return index + 1
}

func (u *Utils) FormatTime(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)
	if duration.Seconds() < 60 {
		return fmt.Sprintf("刚刚")
	} else if duration.Minutes() < 60 {
		return fmt.Sprintf("%.0f 分钟前", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f 小时前", duration.Hours())
	}

	t = t.Add(time.Hour * time.Duration(8))
	return t.Format("2006-01-02 15:04")
}

func (u *Utils) HTML(str string) template.HTML {
	return template.HTML(str)
}

// \n => <br>
func (u *Utils) Br(str string) template.HTML {
	return template.HTML(strings.Replace(str, "\n", "<br>", -1))
}

func (u *Utils) RenderInput(form wtforms.Form, fieldStr string, inputAttrs ...string) template.HTML {
	field, err := form.Field(fieldStr)
	if err != nil {
		panic(err)
	}

	errorClass := ""

	if field.HasErrors() {
		errorClass = " has-error"
	}

	format := `<div class="form-group%s">
        %s
        %s
        %s
    </div>`

	var inputAttrs2 []string = []string{`class="form-control"`}
	inputAttrs2 = append(inputAttrs2, inputAttrs...)

	return template.HTML(
		fmt.Sprintf(format,
			errorClass,
			field.RenderLabel(),
			field.RenderInput(inputAttrs2...),
			field.RenderErrors()))
}

func (u *Utils) RenderInputH(form wtforms.Form, fieldStr string, labelWidth, inputWidth int, inputAttrs ...string) template.HTML {
	field, err := form.Field(fieldStr)
	if err != nil {
		panic(err)
	}

	errorClass := ""

	if field.HasErrors() {
		errorClass = " has-error"
	}
	format := `<div class="form-group%s">
        %s
        <div class="col-lg-%d">
            %s%s
        </div>
    </div>`
	labelClass := fmt.Sprintf(`class="col-lg-%d control-label"`, labelWidth)

	var inputAttrs2 []string = []string{`class="form-control"`}
	inputAttrs2 = append(inputAttrs2, inputAttrs...)

	return template.HTML(
		fmt.Sprintf(format,
			errorClass,
			field.RenderLabel(labelClass),
			inputWidth,
			field.RenderInput(inputAttrs2...),
			field.RenderErrors(),
		))
}

func (u *Utils) HasAd(position string) bool {
	c := DB.C("ads")
	count, _ := c.Find(bson.M{"position": position}).Limit(1).Count()
	return count == 1
}

func (u *Utils) AssertUser(i interface{}) *User {
	v, _ := i.(User)
	return &v
}