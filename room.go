package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
)

type Room struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name: "name"}
}

func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// bindingパッケージでroom生成リクエスト情報をRoomタイプに変換
	r := new(Room)
	errs := binding.Bind(req, r)
	if errs != nil {
		renderer.JSON(w, http.StatusInternalServerError, errs)
		return
	}

	// MongoDBのセッション生成
	session := mongoSession.Copy()
	defer session.Close()

	// MongoDBのID生成
	r.ID = bson.NewObjectId()
	// room情報を保存するため、MongoDBコレクションオブジェクトを生成
	c := session.DB("test").C("rooms")

	// roomsコレクションにroom情報を保存
	if err := c.Insert(r); err != nil {
		// エラー：500
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	// 処理結果を返す
	renderer.JSON(w, http.StatusCreated, r)
}

func retrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	session := mongoSession.Copy()
	defer session.Close()

	var rooms []Room
	// 全てのroom情報を照会
	err := session.DB("test").C("rooms").Find(nil).All(&rooms)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusOK, rooms)
}
