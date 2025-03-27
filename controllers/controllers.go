package controllers

import (
	// "context"
	"fmt"
	// "log"
	"net/http"
	"os"
	// "time"
	// "github.com/coder/websocket"
	// "github.com/coder/websocket/wsjson"
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

// func chatWebSocket(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("we get here")
// 	c, err := websocket.Accept(w, r, nil)
// 	if err != nil {
// 		log.Println("Error Accepting websocket connections", err)
// 	}
// 	defer c.CloseNow()

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()

// 	for {
// 		message := []byte("testing")
// 		var currentIndex, length = 0, len(message)
// 		written, err = c.Write(ctx,,message[currentIndex:length])
// 		if err != nil {
// 			log.Println("Error writing to websocket", err)
// 		}
// 		if currentIndex != length{
// 			currentIndex += written
// 			c.Write(message[currentIndex:length])
// 		}
// 		break
// 	}
// 	var v interface{}
// 	err = wsjson.Read(ctx, c, &v)
// 	if err != nil {
// 		log.Println("Error Reading from websocket", err)
// 	}

// 	log.Printf("received: %v", v)

// 	c.Close(websocket.StatusNormalClosure, "")
// }

func CreateControllers() {
	Router.HandleFunc("/", goConvoHomePage)
	Router.HandleFunc("/chat/{id}", getChatHistory)
	// Router.HandleFunc("/chat/{id}/ws", chatWebSocket) // this will be on the route with "ws" prefix
}
