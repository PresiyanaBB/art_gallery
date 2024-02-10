package main

import (
	"art_gallery/gallery_app"
	"art_gallery/model"
	"art_gallery/mysql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var hasActiveUser bool
var activeUser model.User

type gallery struct {
	app              *gallery_app.ArtGalleryApp
	server           *http.Server
	mux              *http.ServeMux
	templateIndex    *template.Template
	templateLogin    *template.Template
	templateRegister *template.Template
	templateCreate   *template.Template
}

func (g *gallery) Run() error {
	g.mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	g.mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	g.mux.HandleFunc("/", g.handleMain)
	g.mux.HandleFunc("/login", g.handleLogin)
	g.mux.HandleFunc("/register", g.handleRegister)
	g.mux.HandleFunc("/create", g.handleCreate)
	g.templateIndex = template.Must(template.ParseFiles("./templates/index.html"))
	g.templateLogin = template.Must(template.ParseFiles("./templates/login.html"))
	g.templateRegister = template.Must(template.ParseFiles("./templates/register.html"))
	g.templateCreate = template.Must(template.ParseFiles("./templates/create.html"))
	log.Printf("server is listening at %s\n", g.server.Addr)
	if err := g.server.ListenAndServe(); err != nil {
		fmt.Println(fmt.Errorf("failed to start service on port %s:%w", g.server.Addr, err))
		fmt.Print(g.server)
		return nil
	}
	return nil
}

func (g *gallery) handleMain(writer http.ResponseWriter, request *http.Request) {
	g.templateIndex.Execute(writer, hasActiveUser)
}

func (g *gallery) handleLogin(writer http.ResponseWriter, request *http.Request) {
	g.templateLogin.Execute(writer, struct{}{})
}

func (g *gallery) handleRegister(writer http.ResponseWriter, request *http.Request) {
	g.templateRegister.Execute(writer, struct{}{})
}

func (g *gallery) handleCreate(writer http.ResponseWriter, request *http.Request) {
	g.templateCreate.Execute(writer, struct{}{})
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	galleryRepo := mysql.New(mysql.MySQLOptions{
		URI: fmt.Sprintf("%s:%s@tcp(127.0.0.1)/art_gallery",
			"sisi", "sisipresiana03"),
	})
	galleryRepo.Init()
	gallery_appp := gallery_app.New(galleryRepo)
	app := gallery{
		server: server,
		mux:    mux,
		app:    gallery_appp,
	}
	hasActiveUser = false
	app.Run()
}
