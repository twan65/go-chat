package main

import (
	"encoding/json"
	"net/http"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
)

const (
	currentUserKey  = "oauth2_current_use" // セッションに保存されるCurrentUserのキー
	sessionDuration = time.Hour            // ログインのセッション維持時間
)

type User struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"user"`
	AvatarUrl string    `json:"avatar_url"`
	Expired   time.Time `json:"expired"`
}

func (u *User) Valid() bool {
	// 現在時間を基準に満了時間を確認
	return u.Expired.Sub(time.Now()) > 0
}

func (u *User) Refresh() {
	// 満了時間の延長
	u.Expired = time.Now().Add(sessionDuration)
}

func GetCurrentUser(r *http.Request) *User {

	s := sessions.GetSession(r)
	if s.Get(currentUserKey) == nil {
		return nil
	}

	data := s.Get(currentUserKey).([]byte)
	var user User
	json.Unmarshal(data, &user)

	return &user
}

func SetCurrentUser(r *http.Request, user *User) {
	if user != nil {
		user.Refresh()
	}

	s := sessions.GetSession(r)
	val, _ := json.Marshal(user)
	s.Set(currentUserKey, val)
}
