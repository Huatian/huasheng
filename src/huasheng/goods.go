package huasheng

import (
	"labix.org/v2/mgo"
   "labix.org/v2/mgo/bson"
)

type Goods struct {
	Id_  bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Price float32	`bson:"price"`
	Desc string `bson:"desc"`
	Img string `bson:"img"`
}

func loadGoodses() ([]Goods,error){
	var goodses []Goods

	session, db, err := GetDB()
	if err != nil {
		return goodses, err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	
	
	err = db.C("goods").Find(nil).All(&goodses)
	if err != nil {
		return goodses,err
	}
	
	return goodses,nil
}

func loadGoods(name string) (*Goods, error) {
	goods := Goods{}
	session, db, err := GetDB()
	if err != nil {
		return &goods, err
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	err = db.C("goods").Find(bson.M{"name": name}).One(&goods)
	if err != nil {
		return &goods, err
	}

	return &goods, nil
}