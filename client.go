package main

import "golang.org/x/net/websocket"

// 現在、接続中の全てのクライアントリスト
var clients []*Client

type Client struct {
	conn *websocket.Conn // Websocketコネクション
	send chan *Message   // メッセージ送信用チャンネル

	roomId string // 現在接続のチャットルームID
	user   *User  // 現在接続したユーザー情報
}
