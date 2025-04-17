package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"github.com/abhijit360/GoConvo/trace"
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
