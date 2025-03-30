package controllers

import (
	// "context"
	"fmt"
	"time"
	// "log"
	"net/http"
	"os"

	"github.com/abhijit360/GoConvo/sessions"
	// "time"
	"github.com/gorilla/websocket"
)

var Router = http.NewServeMux()
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
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
	
}

func chatWebSocket(w http.ResponseWriter, r *http.Request) {
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
}

func CreateControllers() {
	Router.HandleFunc("/", goConvoHomePage)
	Router.HandleFunc("/chat/{id}", getChatHistory)
	Router.HandleFunc("/create-session",createNewSession)
	// Router.HandleFunc("/chat/{id}/ws", chatWebSocket) // this will be on the route with "ws" prefix
}
