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
	"github.com/coder/websocket"
)

var Router = http.NewServeMux()

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
	// i need to check if thechatroom that is being accessed exists already in the sessionMaanager
	chat_id := r.PathValue("id")
	s, ok := sessions.GetSession(chat_id)
	if !ok{
		var err error
		s, err = sessions.CreateSession(time.Now().String())
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
		}
	}
	// fmt.Println("we get here")
	// c, err := websocket.Accept(w, r, nil)
	// if err != nil {
	// 	log.Println("Error Accepting websocket connections", err)
	// }
	// defer c.CloseNow()

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	// for {
	// 	message := []byte("testing")
	// 	var currentIndex, length = 0, len(message)
	// 	written, err = c.Write(ctx,,message[currentIndex:length])
	// 	if err != nil {
	// 		log.Println("Error writing to websocket", err)
	// 	}
	// 	if currentIndex != length{
	// 		currentIndex += written
	// 		c.Write(message[currentIndex:length])
	// 	}
	// 	break
	// }
	// var v interface{}
	// err = wsjson.Read(ctx, c, &v)
	// if err != nil {
	// 	log.Println("Error Reading from websocket", err)
	// }

	// log.Printf("received: %v", v)

	// c.Close(websocket.StatusNormalClosure, "")
}

func CreateControllers() {
	Router.HandleFunc("/", goConvoHomePage)
	Router.HandleFunc("/chat/{id}", getChatHistory)
	Router.HandleFunc("/create-session",createNewSession)
	// Router.HandleFunc("/chat/{id}/ws", chatWebSocket) // this will be on the route with "ws" prefix
}
