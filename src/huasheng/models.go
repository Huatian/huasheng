package huasheng

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	Id_          bson.ObjectId `bson:"_id"`
	Username     string        `bson:"username"`
	Email        string        `bson:"email"`
	Password     string        `bson:"password"`
	ValidateCode string
	IsSuperuser   bool
	IsActive     bool
	JoinedAt     time.Time
	Index        int
}

// 状态,MongoDB中只存储一个状态
type Status struct {
	Id_        bson.ObjectId `bson:"_id"`
	UserCount  int
	TopicCount int
	ReplyCount int
	UserIndex  int
}
