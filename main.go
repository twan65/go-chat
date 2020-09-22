package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"

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

	// negroniミドルウェア生成
	n := negroni.Classic()

	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	// negroniにrouterをハンドラとして登録
	n.UseHandler(router)

	n.Run(":8000")
}
