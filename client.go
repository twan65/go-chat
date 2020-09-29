package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// 現在、接続中の全てのクライアントリスト
var clients []*Client

const MSG_BUFFER_SIZ = 256

type Client struct {
	conn *websocket.Conn // Websocketコネクション
	send chan *Message   // メッセージ送信用チャンネル

	roomId string // 現在接続のチャットルームID
	user   *User  // 現在接続したユーザー情報
}

func newClient(conn *websocket.Conn, roomId string, u *User) {
	// 新しいクライアントを生成
	c := &Client{
		conn:   conn,
		send:   make(chan *Message, MSG_BUFFER_SIZ),
		roomId: roomId,
		user:   u,
	}

	// clientsに追加
	clients = append(clients, c)

	// メッセージの送・受信を待機
	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) Close() {
	// clientsから終了されたクライアントを削除
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	// sendチャンネルをクローズ
	close(c.send)

	// Websocketのコネクションをクローズ
	c.conn.Close()
	log.Printf("クローズコネクション. アドレス: %s", c.conn.RemoteAddr())
}

func (c *Client) readLoop() {
	// メッセージ受信を待つため、無限ループ
	for {
		m, err := c.read()
		if err != nil {
			log.Println("受信メッセージエラー： ", err)
			break
		}

		// メッセージ受信時、メッセージをDBに生成して全てのclientsに渡す。
		m.create()
		broadcast(m)

	}
	c.Close()
}

func (c *Client) writeLoop() {
	// クライアントのsendチャンネルメッセージの受信待機
	for msg := range c.send {
		// クライアントのチャットルームIDと渡されたメッセージのチャットルームIDが同じの場合Websocketを通してメッセージを渡す。
		if c.roomId == msg.RoomId.Hex() {
			c.write(msg)
		}
	}
}

func broadcast(m *Message) {
	// 全てのクライアントのsendチャンネルにメッセージを渡す。
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) read() (*Message, error) {
	var msg *Message

	// WebsocketコネクションにJSONのメッセージが届いたらMessageタイプでメッセージを読み込む
	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	// Message情報に現在時間とユーザー情報をセット
	msg.CreatedAt = time.Now()
	msg.User = c.user

	log.Println("websocketからのメッセージ:", msg)

	return msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websocket:", m)

	// WebsocketコネクションにJSONでメッセージを渡す。
	return c.conn.WriteJSON(m)
}
