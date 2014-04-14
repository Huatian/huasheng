package huasheng

import (
	"code.google.com/p/go-uuid/uuid"
	"crypto/md5"
	"fmt"
	"github.com/jimmykuu/wtforms"
	"io"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
	"time"
)

// 加密密码,转成md5
func encryptPassword(password string) string {
	h := md5.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 返回当前用户
func currentUser(r *http.Request) (*User, bool) {
	session, _ := store.Get(r, "user")
	username, ok := session.Values["username"]

	if !ok {
		return nil, false
	}

	username = username.(string)

	user := User{}

	c := DB.C("users")

	// 检查用户名
	err := c.Find(bson.M{"username": username}).One(&user)

	if err != nil {
		return nil, false
	}

	return &user, true
}

// URL: /signup
// 处理用户注册,要求输入用户名,密码和邮箱
func signupHandler(w http.ResponseWriter, r *http.Request) {
	form := wtforms.NewForm(
		wtforms.NewTextField("username", "用户名", "", wtforms.Required{}, wtforms.Regexp{Expr: `^[a-zA-Z0-9_]{3,16}$`, Message: "请使用a-z, A-Z, 0-9以及下划线, 长度3-16之间"}),
		wtforms.NewPasswordField("password", "密码", wtforms.Required{}),
		wtforms.NewTextField("email", "电子邮件", "", wtforms.Required{}, wtforms.Email{}),
	)

	if r.Method == "POST" {
		if form.Validate(r) {
			c := DB.C("users")

			result := User{}

			// 检查用户名
			err := c.Find(bson.M{"username": form.Value("username")}).One(&result)
			if err == nil {
				form.AddError("username", "该用户名已经被注册")

				renderTemplate(w, r, "account/signup.html", map[string]interface{}{"form": form})
				return
			}

			// 检查邮箱
			err = c.Find(bson.M{"email": form.Value("email")}).One(&result)

			if err == nil {
				form.AddError("email", "电子邮件地址已经被注册")

				renderTemplate(w, r, "account/signup.html", map[string]interface{}{"form": form})
				return
			}

			c2 := DB.C("status")
			var status Status
			c2.Find(nil).One(&status)

			id := bson.NewObjectId()
			username := form.Value("username")
			validateCode := strings.Replace(uuid.NewUUID().String(), "-", "", -1)
			index := status.UserIndex + 1
			err = c.Insert(&User{
				Id_:          id,
				Username:     username,
				Password:     encryptPassword(form.Value("password")),
				Email:        form.Value("email"),
				ValidateCode: validateCode,
				IsActive:     true,
				JoinedAt:     time.Now(),
				Index:        index,
			})

			if err != nil {
				panic(err)
			}

			c2.Update(nil, bson.M{"$inc": bson.M{"userindex": 1, "usercount": 1}})

			// 发送邮件
			/*
							subject := "欢迎加入Golang 中国"
							message2 := `欢迎加入Golang 中国。请访问下面地址激活你的帐户。

				<a href="%s/activate/%s">%s/activate/%s</a>

				如果你没有注册，请忽略这封邮件。

				©2012 Golang 中国`
							message2 = fmt.Sprintf(message2, config["host"], validateCode, config["host"], validateCode)
							sendMail(subject, message2, []string{formstore.Value("email")})

							message(w, r, "注册成功", "请查看你的邮箱进行验证，如果收件箱没有，请查看垃圾邮件，如果还没有，请给jimmykuu@126.com发邮件，告知你的用户名。", "success")
			*/
			// 注册成功后设成登录状态
			session, _ := store.Get(r, "user")
			session.Values["username"] = username
			session.Save(r, w)

			// 跳到修改用户信息页面
			http.Redirect(w, r, "/profile", http.StatusFound)
			return
		}
	}

	renderTemplate(w, r, "account/signup.html", map[string]interface{}{"form": form})
}

// URL: /signin
// 处理用户登录,如果登录成功,设置Cookie
func signinHandler(w http.ResponseWriter, r *http.Request) {
	next := r.FormValue("next")

	form := wtforms.NewForm(
		wtforms.NewHiddenField("next", next),
		wtforms.NewTextField("username", "用户名", "", &wtforms.Required{}),
		wtforms.NewPasswordField("password", "密码", &wtforms.Required{}),
	)

	if r.Method == "POST" {
		if form.Validate(r) {
			c := DB.C("users")
			user := User{}

			err := c.Find(bson.M{"username": form.Value("username")}).One(&user)

			if err != nil {
				form.AddError("username", "该用户不存在")

				renderTemplate(w, r, "account/signin.html", map[string]interface{}{"form": form})
				return
			}

			if !user.IsActive {
				form.AddError("username", "邮箱没有经过验证,如果没有收到邮件,请联系管理员")
				renderTemplate(w, r, "account/signin.html", map[string]interface{}{"form": form})
				return
			}

			if user.Password != encryptPassword(form.Value("password")) {
				form.AddError("password", "密码和用户名不匹配")

				renderTemplate(w, r, "account/signin.html", map[string]interface{}{"form": form})
				return
			}

			session, _ := store.Get(r, "user")
			session.Values["username"] = user.Username
			session.Save(r, w)

			if form.Value("next") == "" {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				http.Redirect(w, r, next, http.StatusFound)
			}

			return
		}
	}

	renderTemplate(w, r, "account/signin.html", map[string]interface{}{"form": form})
}