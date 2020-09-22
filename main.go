package main

import (
	"net/http"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/urfave/negroni"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var renderer *render.Render

const (
	sessionKey    = "twan_chat_session"
	sessionSecret = "twan_chat_session_secret"
)

func init() {
	// renderer生成
	renderer = render.New()
}

func main() {

	// Router生成
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Twan Chat!"})
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

	// negroniミドルウェア生成
	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	// negroniにrouterをハンドラとして登録
	n.UseHandler(router)

	n.Run(":8000")
}
