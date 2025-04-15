package main

import "github.com/gorilla/websocket"

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type chatRoom struct {
	// this channel holds the incoming message
	// that must be forwarded to other clients
	triage chan []byte

	// a channel for clients to join the chatRoom
	join chan *client

	// a channel for clients to leave the chatroom
	leave chan *client

	clients map[*client]bool
}

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize, WriteBufferSize: socketBufferSize}

func (cr *chatRoom) ServeHTTP(w http.ResponseWriter, req *http.Request){
	socket, err := upgrader.Upgrade(w,req,nil)
	if err != nil {
		log.Fatal("Failed to upgrade to websocket connection", err)
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		chatRoom: cr,
	}
	cr.join <- client

	defer func() {cr.leave <- client}()
	go client.write()
	client.read()
}

func (chatRoom *chatRoom) run() {
	for {
		select {
		case client := <-chatRoom.join:
			// adding client
			chatRoom.clients[client] = true
		case client := <-chatRoom.leave:
			// rmeove client
			chatRoom.clients[client] = false
		case msg := <-chatRoom.triage:
			for client := range chatRoom.clients {
				client.send <- msg
			}
		}
	}
}

func newRoom() *chatRoom {
	return &chatRoom{
		triage: make(chan []byte),
		join: make(chan * client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}