package main

import (
	"art_gallery/gallery_app"
	model "art_gallery/models"
	"art_gallery/mysql"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	//load static files - images and css
	g.mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	g.mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	//handle functions
	g.mux.HandleFunc("/", g.handleMain)
	g.mux.HandleFunc("/login", g.handleLogin)
	g.mux.HandleFunc("/logout", g.handleLogout)
	g.mux.HandleFunc("/register", g.handleRegister)
	g.mux.HandleFunc("/create", g.handleCreate)
	g.mux.HandleFunc("/enter_account", g.handleEnterAccount)
	g.mux.HandleFunc("/create_account", g.handleCreateAccount)
	// g.mux.HandleFunc("/", g.handle)
	//load html templates
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
	paintings, err := g.app.GetAll()
	if err != nil {
		log.Printf("failed to get posts: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := struct {
		Active bool
		User   model.User
		Paint  []model.Painting
	}{
		Active: hasActiveUser,
		User:   activeUser,
		Paint:  paintings,
	}

	writer.WriteHeader(http.StatusOK)
	g.templateIndex.Execute(writer, data)
}

func (g *gallery) handleLogin(writer http.ResponseWriter, request *http.Request) {
	g.templateLogin.Execute(writer, struct{}{})
}

func (g *gallery) handleLogout(writer http.ResponseWriter, request *http.Request) {
	hasActiveUser = false
	http.Redirect(writer, request, "/", http.StatusFound)
}

func (g *gallery) handleRegister(writer http.ResponseWriter, request *http.Request) {
	g.templateRegister.Execute(writer, struct{}{})
}

func (g *gallery) handleCreate(writer http.ResponseWriter, request *http.Request) {
	g.templateCreate.Execute(writer, struct{}{})
}

func (g *gallery) handleEnterAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	email := request.Form.Get("email")
	password := request.Form.Get("password")

	hash := sha256.New()
	hash.Write([]byte(password))
	hashedBytes := hash.Sum(nil)
	hashedPassword := hex.EncodeToString(hashedBytes)

	user, _ := g.app.FindByUserEmail(email)
	if err != nil || user == nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	if user.Password == hashedPassword {
		hasActiveUser = true
		activeUser = *user
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func (g *gallery) handleCreateAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	email := request.Form.Get("email")
	fname := request.Form.Get("first-name")
	lname := request.Form.Get("last-name")
	password := request.Form.Get("password")
	re_password := request.Form.Get("re-password")

	if password != re_password {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(password))
	hashedBytes := hash.Sum(nil)
	hashedPassword := hex.EncodeToString(hashedBytes)

	user, _ := g.app.FindByUserEmail(email)
	if err != nil {
		http.Redirect(writer, request, "/register", http.StatusFound)
		return
	}

	if user != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	var new_user model.User
	new_user.ID = uuid.New().String()
	new_user.DateOfRegistration = time.Now()
	new_user.Email = email
	new_user.Password = hashedPassword
	new_user.FirstName = fname
	new_user.LastName = lname

	hasActiveUser = true
	activeUser = new_user
	g.app.AddUser(&activeUser)
	http.Redirect(writer, request, "/", http.StatusFound)
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
