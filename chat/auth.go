package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie{
		// Not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}

	// call the next hander
	h.next.ServeHTTP(w,r)
}

func RequireAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

/* 
this login handler is function and not of type Handler because it does not need to store any state!
*/
func loginHandler(w http.ResponseWriter, r *http.Request){
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	
	switch action {
	case "login":
		log.Println("TODO handle login for",provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported",action)
	}
}