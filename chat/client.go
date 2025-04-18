package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// c;oent represents a single chatting user
type client struct {
	// the websocket connection for this client
	socket *websocket.Conn
	
	// this is the channel on which messsages are going to be sent 
	send chan *message
	
	// room is the room this client is chatting in
	chatRoom *chatRoom

	// userData holds information about the user
	userData map[string]interface{}
}

func (c *client) read(){
	defer c.socket.Close()
	for {
		var msg message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		c.chatRoom.triage <- &msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg:= range c.send{
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return 
		}
	}
}