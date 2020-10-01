package main

import (
	"log"
	"net/http"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

const (
	sessionKey         = "twan_chat_session"
	sessionSecret      = "twan_chat_session_secret"
	SOCKET_BUFFER_SIZE = 1024
)

var (
	renderer     *render.Render
	mongoSession *mgo.Session
	upgrader     = &websocket.Upgrader{
		ReadBufferSize:  SOCKET_BUFFER_SIZE,
		WriteBufferSize: SOCKET_BUFFER_SIZE,
	}
)

func init() {
	// renderer生成
	renderer = render.New()

	// MongoDBの情報（認証情報を含む）
	mongoInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Timeout:  20 * time.Second,
		Database: "",
		Username: "twan",
		Password: "1234",
		Source:   "",
	}

	// MongoDB接続セッションを作成
	s, err := mgo.DialWithInfo(mongoInfo)
	if err != nil {
		panic(err)
	}

	s.SetMode(mgo.Monotonic, true)
	mongoSession = s
}

func main() {

	// Router生成
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]interface{}{"host": r.Host})
	})

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// ログイン画面を表示
		renderer.HTML(w, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// セッション情報を削除し、ログイン画面に遷移
		sessions.GetSession(r).Delete(currentUserKey)
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)
	router.GET("/rooms", getRooms)
	router.GET("/rooms/:id", getRoom)
	router.POST("/rooms", createRoom)
	router.DELETE("/rooms/:id", deleteRoom)
	router.DELETE("/messages/:room_id/:id", deleteMessage)

	router.GET("/rooms/:id/messages", getMessages)

	router.GET("/ws/:room_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		newClient(socket, ps.ByName("room_id"), GetCurrentUser(r))
	})

	// negroniミドルウェア生成
	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	// LoginRequiredハンドラをnegroniに登録
	n.Use(LoginRequired("/login", "/auth"))

	// negroniにrouterをハンドラとして登録
	n.UseHandler(router)

	n.Run(":8000")
}
