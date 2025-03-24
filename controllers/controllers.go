package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var Router = http.NewServeMux()

func goConvoHomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("getting the home page")
	}
}

func getChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Printf("getting the chat with id %v", r.PathValue("id"))
	}
}

func chatWebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Error Accepting websocket connections", err)
	}
	defer c.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		log.Println("Error Reading from websocket", err)
	}

	log.Printf("received: %v", v)

	c.Close(websocket.StatusNormalClosure, "")
}

func CreateControllers() {
	Router.HandleFunc("/", goConvoHomePage)
	Router.HandleFunc("/chat/{id}", getChatHistory)
	Router.HandleFunc("/chat/{id}/", chatWebSocket) // this will be on the route with "ws" prefix
}
