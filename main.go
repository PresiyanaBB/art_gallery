package main

import (
	"art_gallery/gallery_app"
	model "art_gallery/models"
	"art_gallery/mysql"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	templateAccount  *template.Template
	templateEdit     *template.Template
}

func (g *gallery) Run() error {
	// for _, v := range model.GenreTypesString {
	// 	var genre model.Genre
	// 	genre.ID = uuid.New().String()
	// 	genre.Name = v
	// 	g.app.AddGenre(&genre)
	// }

	//load static files - images and css
	g.mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	g.mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	//handle functions
	//user auth
	g.mux.HandleFunc("/", g.handleMain)
	g.mux.HandleFunc("/login", g.handleLogin)
	g.mux.HandleFunc("/logout", g.handleLogout)
	g.mux.HandleFunc("/register", g.handleRegister)
	g.mux.HandleFunc("/enter_account", g.handleEnterAccount)
	g.mux.HandleFunc("/create_account", g.handleCreateAccount)
	g.mux.HandleFunc("/account", g.handleAccount)
	//painting manager
	g.mux.HandleFunc("/create", g.handleCreate)
	g.mux.HandleFunc("/create_painting", g.handleCreatePainting)
	g.mux.HandleFunc("/edit_painting", g.handleEditPainting)
	g.mux.HandleFunc("/update_editted_painting", g.handleUpdateEdittedPainting)
	g.mux.HandleFunc("/delete_painting", g.handleDeletePainting)
	//load html templates
	g.templateIndex = template.Must(template.ParseFiles("./templates/index.html"))
	g.templateLogin = template.Must(template.ParseFiles("./templates/login.html"))
	g.templateRegister = template.Must(template.ParseFiles("./templates/register.html"))
	g.templateCreate = template.Must(template.ParseFiles("./templates/create.html"))
	g.templateAccount = template.Must(template.ParseFiles("./templates/account.html"))
	g.templateEdit = template.Must(template.ParseFiles("./templates/edit.html"))

	log.Printf("server is listening at %s\n", g.server.Addr)
	if err := g.server.ListenAndServe(); err != nil {
		fmt.Println(fmt.Errorf("failed to start service on port %s:%w", g.server.Addr, err))
		fmt.Print(g.server)
		return nil
	}
	return nil
}

func (g *gallery) handleMain(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	paintings, err := g.app.GetAll()
	genres, _ := g.app.GetAllGenres()

	//parse form
	sb_author := request.Form.Get("sb_author")
	sb_title := request.Form.Get("sb_title")
	category := request.Form.Get("sb_category")
	if sb_author == "" && sb_title == "" && category == "" {
	} else if sb_author != "" && sb_title != "" && category != "" {
		paintings, err = g.app.FindByUserNameAndPaintingTitleAndCenre(sb_author, sb_title, category)
	} else if sb_author == "" && sb_title != "" && category != "" {
		paintings, err = g.app.FindByPaintingTitleAndGenre(sb_title, category)
	} else if sb_author != "" && sb_title == "" && category != "" {
		paintings, err = g.app.FindByUserNameAndGenre(sb_author, category)
	} else if sb_author != "" && sb_title != "" && category == "" {
		paintings, err = g.app.FindByUserNameAndPaintingTitle(sb_author, sb_title)
	} else if sb_author != "" && sb_title == "" && category == "" {
		paintings, err = g.app.FindByUserName(sb_author)
	} else if sb_title != "" && sb_author == "" && category == "" {
		paintings, err = g.app.FindByTitle(sb_title)
	} else if category != "" && sb_title == "" && sb_author == "" {
		paintings, err = g.app.FindByGenre(category)
	}

	if err != nil {
		log.Printf("failed to get posts: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(paintings) == 0 {
		paintings, _ = g.app.GetAll()
	}

	var encodedData []string

	for _, v := range paintings {
		encodedData = append(encodedData, base64.StdEncoding.EncodeToString(v.Data))
	}

	data := struct {
		Active        bool
		User          model.User
		Paint         []model.Painting
		GenreTypes    []string
		EncodedImages []string
	}{
		Active:        hasActiveUser,
		User:          activeUser,
		Paint:         paintings,
		GenreTypes:    genres,
		EncodedImages: encodedData,
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
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
	if err != nil || user == nil || hashedPassword != user.Password {
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

func (g *gallery) handleAccount(writer http.ResponseWriter, request *http.Request) {
	paintings, err := g.app.GetAll()
	if err != nil {
		log.Printf("failed to get posts: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var encodedData []string

	for _, v := range paintings {
		encodedData = append(encodedData, base64.StdEncoding.EncodeToString(v.Data))
	}

	data := struct {
		User          model.User
		Paint         []model.Painting
		EncodedImages []string
	}{
		User:          activeUser,
		Paint:         paintings,
		EncodedImages: encodedData,
	}

	writer.WriteHeader(http.StatusOK)
	g.templateAccount.Execute(writer, data)
}

func (g *gallery) handleCreate(writer http.ResponseWriter, request *http.Request) {
	genres, _ := g.app.GetAllGenres()
	g.templateCreate.Execute(writer, genres)
}

func (g *gallery) handleCreatePainting(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	////////

	file, handler, err := request.FormFile("imageUrl")
	if err != nil {
		http.Error(writer, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file in the server to store the uploaded photo temporarily
	tempFile, err := os.CreateTemp("", "temp-*.png")
	if err != nil {
		http.Error(writer, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(writer, "Error copying file", http.StatusInternalServerError)
		return
	}

	// Read the temporary file content
	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		http.Error(writer, "Error reading file", http.StatusInternalServerError)
		return
	}

	name := request.Form.Get("name")
	category := request.Form.Get("category")
	height := request.Form.Get("height")
	width := request.Form.Get("width")
	description := request.Form.Get("description")
	price := request.Form.Get("price")

	var painting model.Painting
	painting.ID = uuid.New().String()
	painting.Author = activeUser
	painting.DateOfPublication = time.Now()
	painting.Description = description
	painting.Price, _ = strconv.ParseFloat(price, 64)
	painting.Title = name
	painting.MIMEType = handler.Header.Get("Content-Type")
	painting.Data = data
	var gg *model.Genre
	gg, _ = g.app.FindGenre(category)
	if gg != nil {
		painting.Genre = *gg
	}
	painting.Height, _ = strconv.Atoi(height)
	painting.Width, _ = strconv.Atoi(width)

	g.app.AddPainting(&painting)
	http.Redirect(writer, request, "/account", http.StatusFound)
}

func (g *gallery) handleDeletePainting(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	painting_id := request.Form.Get("painting_id")
	g.app.DeletePaintingByID(painting_id)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func (g *gallery) handleEditPainting(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	painting_id := request.Form.Get("painting_id")
	p, _ := g.app.FindPaintingByID(painting_id)
	genres, _ := g.app.GetAllGenres()

	data := struct {
		Painting   model.Painting
		GenreTypes []string
	}{
		Painting:   *p,
		GenreTypes: genres,
	}

	writer.WriteHeader(http.StatusOK)
	g.templateEdit.Execute(writer, data)
}

func (g *gallery) handleUpdateEdittedPainting(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Printf("error parsing html form: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	painting_id := request.Form.Get("painting_id")
	currentPainting, _ := g.app.FindPaintingByID(painting_id)

	name := request.Form.Get("name")
	category := request.Form.Get("category")
	height := request.Form.Get("height")
	width := request.Form.Get("width")
	description := request.Form.Get("description")
	price := request.Form.Get("price")

	var painting model.Painting
	painting.ID = uuid.New().String()
	painting.Author = activeUser
	painting.DateOfPublication = time.Now()
	painting.Description = description
	painting.Price, _ = strconv.ParseFloat(price, 64)
	painting.Title = name
	painting.MIMEType = currentPainting.MIMEType
	painting.Data = currentPainting.Data
	var gg *model.Genre
	gg, _ = g.app.FindGenre(category)
	if gg != nil {
		painting.Genre = *gg
	}
	painting.Height, _ = strconv.Atoi(height)
	painting.Width, _ = strconv.Atoi(width)

	g.app.DeletePaintingByID(painting_id)
	g.app.AddPainting(&painting)
	writer.Header().Set("Content-Type", "text/html")
	http.Redirect(writer, request, "/account", http.StatusFound)
}

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	galleryRepo := mysql.New(mysql.MySQLOptions{
		URI: fmt.Sprintf("%s:%s@tcp(127.0.0.1)/art_gallery",
			"{name}", "{password}"),
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
