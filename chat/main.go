package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/abhijit360/GoConvo/trace"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"

	// "github.com/stretchr/gomniauth/providers/facebook"
	// "github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// loads the source file, compiles the template and executes it
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error retrieving auth key or auth secret from env file. .env file is potentially missing")
	}
	// setup goOmniAuth
	gomniauth.SetSecurityKey("my-auth-key")
	gomniauth.WithProviders(
		// facebook.New("key", "secret",
		// 	"http://localhost:8080/auth/callback/facebook"),
		// github.New("key", "secret",
		// 	"http://localhost:8080/auth/callback/github"),
		google.New(os.Getenv("GOOGLE_AUTH_KEY"), os.Getenv("GOOGLE_AUTH_SECRET"),
			"http://localhost:8080/auth/callback/google"),
	)

	chat := newRoom()
	chat.tracer = trace.New(os.Stdout)
	http.Handle("/chat", RequireAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/chatRoom", chat)
	go chat.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
