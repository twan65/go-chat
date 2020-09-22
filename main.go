package main

import (
	"net/http"

	"github.com/codegangsta/negroni"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var renderer *render.Render

func init() {
	// renderer生成
	renderer = render.New()
}

func main() {

	// Router生成
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple Chat!"})
	})

	// negroniミドルウェア生成
	n := negroni.Classic()

	// negroniにrouterをハンドラとして登録
	n.UseHandler(router)

	n.Run(":8000")
}
