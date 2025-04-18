package main

import (
	"log"
	"net/http"

	"github.com/abhijit360/GoConvo/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type chatRoom struct {
	// this channel holds the incoming message
	// that must be forwarded to other clients
	triage chan *message

	// a channel for clients to join the chatRoom
	join chan *client

	// a channel for clients to leave the chatroom
	leave chan *client

	clients map[*client]bool

	tracer trace.Tracer
}

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize, WriteBufferSize: socketBufferSize}

func (cr *chatRoom) ServeHTTP(w http.ResponseWriter, req *http.Request){
	socket, err := upgrader.Upgrade(w,req,nil)
	if err != nil {
		log.Fatal("Failed to upgrade to websocket connection", err)
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie when upgrading connection:", err)
		return
	}
	client := &client{
		socket: socket,
		send: make(chan *message, messageBufferSize),
		chatRoom: cr,
		userData: objx.MustFromBase64(authCookie.Value),
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
			chatRoom.tracer.Trace("New Client joined")
		case client := <-chatRoom.leave:
			// rmeove client
			chatRoom.clients[client] = false
			chatRoom.tracer.Trace("Client left")
		case msg := <-chatRoom.triage:
			for client := range chatRoom.clients {
				client.send <- msg
				chatRoom.tracer.Trace("message recieved: ",msg.Message)
			}
		}
	}
}

func newRoom() *chatRoom {
	return &chatRoom{
		triage: make(chan *message),
		join: make(chan * client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}