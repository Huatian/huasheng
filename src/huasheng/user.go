package huasheng

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Name     string        `bson:"name"`
	Email    string        `bson:"email"`
	Password string        `bson:"password"`
}

//用于用户注册，一旦保存，不可更改
func (u *User) Create() error {
	session, db, err := GetDB()
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return db.C("user").Insert(u)
}

//用于用户登录验证
func (u *User) Login() (bool, error) {
	session, db, err := GetDB()
	if err != nil {
		return false, err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	user := User{}
	err = db.C("user").Find(bson.M{"name": u.Name}).One(&user)
	if err != nil {
		return false,err
	}
	
	return u.Password == user.Password,nil
}

func (u *User) IsExist() (bool, error) {
	session, db, err := GetDB()
	if err != nil {
		return false, err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	user := User{}
	err = db.C("user").Find(bson.M{"name": u.Name}).One(&user)
	if err != nil {
		return false,err
	}
	
	return u.Name == user.Name,nil
}

type UserInfo struct {
	Id_     bson.ObjectId `bson:"_id"`
	Gender  bool          `bson:"gender"`
	Age     int           `bson:"age"`
	Address string        `bson:"address"`
}

func (u *UserInfo) Save() error {
	session, db, err := GetDB()
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return db.C("userinfo").Insert(u)
}

func (u *UserInfo) Update() error {
	session, db, err := GetDB()
	if err != nil {
		return err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return db.C("userinfo").Update(bson.M{"_id": u.Id_}, bson.M{"$set": bson.M{"gender": u.Gender, "age": u.Age, "address": u.Address}})
}
