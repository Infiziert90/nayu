package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
)

type Upload struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	UniqueCode string        `bson:"ucode" json:"ucode"`
	DeleteCode string        `bson:"dcode" json:"dcode"`
	File       string        `bson:"file" json:"file"`
}

type UploadDAO struct {
	Server   string
	Database string
}

const COLLECTION = "uploads"

var (
	db     *mgo.Database
	Dao    = UploadDAO{"localhost", "filehoster"}
	UsedUC = make(map[string]bool)
)

func (u *UploadDAO) Connect() {
	session, err := mgo.Dial(u.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(u.Database)
	u.createUsedUC()
}

func (u *UploadDAO) createUsedUC() {
	var upload []Upload
	db.C(COLLECTION).Find(bson.M{}).All(&upload)
	for x := range upload {
		UsedUC[upload[x].UniqueCode] = true
	}
}

func (u *UploadDAO) FindByUC(uc string) (Upload, error) {
	var upload Upload
	err := db.C(COLLECTION).Find(bson.M{"ucode": uc}).One(&upload)
	return upload, err
}

func (u *UploadDAO) FindByDC(dc string) (Upload, error) {
	var upload Upload
	err := db.C(COLLECTION).Find(bson.M{"dcode": dc}).One(&upload)
	return upload, err
}

func (u *UploadDAO) Insert(upload *Upload) error {
	err := db.C(COLLECTION).Insert(&upload)
	return err
}

func (u *UploadDAO) Delete(upload *Upload) error {
	err := db.C(COLLECTION).Remove(&upload)
	return err
}
