package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

const MSG_FETCH_SIZE = 10

type Message struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	RoomId    bson.ObjectId `bson:"room_id" json:"room_id"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	User      *User         `bson:"user" json:"user"`
}

func (m *Message) create() error {
	session := mongoSession.Copy()
	defer session.Close()

	// Create MongDB ID
	m.ID = bson.NewObjectId()

	// メッセージ生成時間を記録
	m.CreatedAt = time.Now()

	c := session.DB("test").C("messages")

	// messagesコレクションにmessage情報を保存
	if err := c.Insert(m); err != nil {
		return err
	}
	return nil
}

func retrieveMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// MongDBのセッション生成
	session := mongoSession.Copy()
	defer session.Close()

	// クエリパラメータのlimit値の確認
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		// 正常のlimit値ではない場合limitをmessageFetchSizeでセット
		limit = MSG_FETCH_SIZE
	}

	var messages []Message
	// _idを逆順でソートし、limitの数分message照会
	err = session.DB("test").C("messages").
		Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).Sort("-_id").Limit(limit).All(&messages)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	// 結果を返す
	renderer.JSON(w, http.StatusOK, messages)
}

// TODO
func deleteMessage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	colQuerier := bson.M{"room_id": bson.ObjectIdHex(ps.ByName("room_id")), "_id": bson.ObjectIdHex(ps.ByName("id"))}

	err := session.DB("test").C("messages").RemoveId(colQuerier)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusNoContent, nil)
}
