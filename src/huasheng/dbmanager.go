package huasheng

import (
	"fmt"
	"labix.org/v2/mgo"
)


func GetDB() (*mgo.Session, *mgo.Database, error) {
	session, err1 := mgo.Dial("localhost")
	if err1 != nil {
		fmt.Println("Session error:", err1)
		return session, nil, err1
	}
	db := session.DB("huasheng")
	return session, db, nil
}

