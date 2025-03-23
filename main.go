package main

import (
	"fmt"
	"net/http"
	"github.com/abhijit360/GoConvo/controllers"
)


func main() {
	fmt.Println("starting program")
	controllers.CreateControllers()
	server := &http.Server{
		Addr: ":8080",
		Handler: controllers.Router,
	}
	server.ListenAndServe()
}