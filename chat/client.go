package main

import (
	"github.com/gorilla/websocket"
)

// c;oent represents a single chatting user
type client struct {
	// the websocket connection for this client
	socket *websocket.Conn
	
	// this is the channel on which messsages are going to be sent 
	send chan []byte
	
	// room is the room this client is chatting in
	chatRoom *chatRoom
}

func (c *client) read(){
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.chatRoom.triage <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg:= range c.send{
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return 
		}
	}
}