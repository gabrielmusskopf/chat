package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var port = ":9000";
var server = NewServer()

type Server struct {
	conns map[*websocket.Conn]bool
	chats []*Chat
}

type Client struct {
	conn *websocket.Conn
}

type Chat struct {
	Id      int
	Name    string
	Clients []*Client
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		conn: ws,
	}
}

func NewChat(id int, name string) *Chat {
	return &Chat{
		Id:   id,
		Name: name,
	}
}

func (c *Chat) addClient(client *Client) {
	c.Clients = append(c.Clients, client)
}

func (s *Server) addChat(c *Chat) {
	s.chats = append(s.chats, c)
}

func (s *Server) GetChat(id int) *Chat {
	for _, chat := range s.chats {
		if chat.Id == id {
			return chat
		}
	}
	return nil
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

func (c *Chat) readLoop(ws *websocket.Conn) {
	for {
		msgType, buf, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("unexpected close error:", err)
				break
			}
			continue
		}

		c.send(msgType, buf)
	}
}

func (c *Chat) send(msgType int, b []byte) {
	for _, client := range c.Clients {
		go func(ws *websocket.Conn) {
			if err := ws.WriteMessage(msgType, b); err != nil {
				fmt.Println("write error", err)
			}
		}(client.conn)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleChatWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("id invalid")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("ws upgrade err:", err)
		return
	}

	client := NewClient(conn)
	chat := server.GetChat(id)

	if chat == nil {
		chat = NewChat(id, "chat")
		server.addChat(chat)

		fmt.Println("chat", id, "created")
	}

	chat.addClient(client)
	chat.send(websocket.TextMessage, []byte("client entered chat!"))
    fmt.Println("client entered chat ", chat.Id)

	chat.readLoop(conn)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws/chat/{id}", handleChatWS)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.Handle("/", r)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
