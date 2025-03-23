package controllers

import (
	"fmt"
	"net/http"
)

var Router = http.NewServeMux()

func goConvoHomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("getting the home page")
	}
}

func getChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("getting the chat with id %v", r.PathValue("id"))
	}
}

func CreateControllers() {
	Router.HandleFunc("/", goConvoHomePage)
	Router.HandleFunc("/chat/{id}", getChatHistory)
}
