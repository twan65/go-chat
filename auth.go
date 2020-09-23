package main

import (
	"log"
	"net/http"
	"strings"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/urfave/negroni"
)

const (
	nextPageKey     = "next_page"
	authSecurityKey = "auth_security_key"
)

func init() {
	// gomniauth情報セット
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(
		google.New(google_client_id, google_client_secret, "http://localhost:8000/auth/callback/google"),
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	s := sessions.GetSession(r)

	switch action {
	case "login":
		// gomniauth.Providerのlogin画面に遷移
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, r, loginUrl, http.StatusFound)
	case "callback":
		// gomniauth callback処理
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}

		// コールバック結果からユーザー情報確認
		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		if err != nil {
			log.Fatalln(err)
		}

		u := &User{
			Uid:       user.Data().Get("id").MustStr(),
			Name:      user.Name(),
			Email:     user.Email(),
			AvatarUrl: user.AvatarURL(),
		}

		SetCurrentUser(r, u)
		http.Redirect(w, r, s.Get(nextPageKey).(string), http.StatusFound)
	default:
		http.Error(w, "Auth action '"+action+"' is not supported", http.StatusNotFound)
	}

}

func LoginRequired(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// ignore urlの場合は次のハンドラを実行
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}
		// CurrentUser情報取得
		u := GetCurrentUser(r)

		// CurrentUser情報が有効の場合、満了時間を更新して次のハンドラを実行
		if u != nil && u.Valid() {
			SetCurrentUser(r, u)
			next(w, r)
			return
		}

		// CurrentUser情報が有効ではない場合はCurrentUserをnilでセット
		SetCurrentUser(r, nil)

		// ログイン後、遷移するURLをセッションに保存(r)
		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())

		// ログイン画面にリダイレクト
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
