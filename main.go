package main

import (
	"fmt"
	"github.com/abhijit360/GoConvo/controllers"
	"net/http"
)

func main() {
	fmt.Println("starting program")
	controllers.CreateControllers()
	server := &http.Server{
		Addr:    ":8080",
		Handler: controllers.Router,
	}
	server.ListenAndServe()
}
