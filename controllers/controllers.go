package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abhijit360/GoConvo/sessions"
	"github.com/gorilla/websocket"
)

var Router = http.NewServeMux()
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type createSessionRequest struct {
	currentTime string `json:currentTime`
}

func goConvoHomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" || r.URL.Path != "/" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	fmt.Println("we are getting this far")
	dir, _ := os.Getwd()
	fmt.Println("current dir", dir)
	http.ServeFile(w, r, "./templates/index.html")
}

func getChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "not found", http.StatusNotFound)
	}

	fmt.Printf("getting the chat with id %v", r.PathValue("id"))
}

func createNewSession(w http.ResponseWriter, r * http.Request){
	if r.Method != "POST"{
		http.Error(w,"Wrong HTTP method to create session",http.StatusBadRequest)
	}
	var response createSessionRequest
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		fmt.Errorf("unable to parse json %v",err)
		return
	}
	s, err := sessions.CreateSession(response.currentTime)

	go s.HandleBroadcast() // create the broadcast

	if err != nil {
		fmt.Errorf("unable to create session %v",err)
	}
	chatMetaData := s.GetChatMetaData()
	w.Write([]byte(chatMetaData.Chat_id))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// create session or get existing session
	chat_id := r.PathValue("id")
	s, ok := sessions.GetSession(chat_id)
	if !ok{
		var err error
		s, err = sessions.CreateSession(time.Now().String())
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
		}
	}
	
	// intercept websocket connection
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		fmt.Errorf("Having trouble upgrading connection",err)
		return
	}
	s.AddSession(conn)

	defer conn.Close()
	defer s.RemoveSession(conn)
	done := make(chan struct{})

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				close(done)
				break
			}
			s.Broadcast <- message
		}
	}()

	/* 
	this maintains the connection until we want to explicitly break out.
	if we do not have this we could end up with the race condition
	where we end the connection before we even read from the connection 
	due to the concurrent aspect of the go routine */
	<- done
}

func CreateControllers() {
	Router.HandleFunc("/", goConvoHomePage)
	Router.HandleFunc("/chat/{id}", getChatHistory)
	Router.HandleFunc("/create-session",createNewSession)
	// Router.HandleFunc("/chat/{id}/ws", chatWebSocket) // this will be on the route with "ws" prefix
}
