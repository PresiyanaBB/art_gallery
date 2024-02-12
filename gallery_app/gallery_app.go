package gallery_app

import model "art_gallery/models"

type PaintingRepository interface {
	GetAll() ([]model.Painting, error)
	AddPainting(painting *model.Painting) error
	DeleteAll() error
	FindByTitle(title string) ([]model.Painting, error)
	FindBySize(width int, height int) ([]model.Painting, error)
	FindByGenre(genre string) ([]model.Painting, error)
	FindByUserName(name string) ([]model.Painting, error)
	FindByUserEmail(email string) (*model.User, error)
	AddUser(u *model.User) error
	GetAllGenres() ([]string, error)
	FindGenre(genre string) (*model.Genre, error)
	AddGenre(genre *model.Genre) error
	DeleteAllGenres() error
	FindGenreByID(id string) (*model.Genre, error)
	FindUserByID(id string) (*model.User, error)
	FindPaintingByID(id string) (*model.Painting, error)
	DeletePaintingByID(id string) error
	FindGenreByName(name string) (*model.Genre, error)
	FindUsersByFirstName(name string) ([]model.User, error)
	FindByUserNameAndPaintingTitle(name string, title string) ([]model.Painting, error)
	FindByUserNameAndGenre(name string, genreName string) ([]model.Painting, error)
	FindByPaintingTitleAndGenre(title string, genreName string) ([]model.Painting, error)
	FindByUserNameAndPaintingTitleAndCenre(name string, title string, genreName string) ([]model.Painting, error)
}

type ArtGalleryApp struct {
	paintings PaintingRepository
}

func New(painting PaintingRepository) *ArtGalleryApp {
	return &ArtGalleryApp{
		paintings: painting,
	}
}

func (g *ArtGalleryApp) GetAll() ([]model.Painting, error) {
	return g.paintings.GetAll()
}

func (g *ArtGalleryApp) AddPainting(p *model.Painting) error {
	return g.paintings.AddPainting(p)
}

func (g *ArtGalleryApp) DeleteAll() error {
	return g.paintings.DeleteAll()
}

func (g *ArtGalleryApp) FindByTitle(title string) ([]model.Painting, error) {
	return g.paintings.FindByTitle(title)
}

func (g *ArtGalleryApp) FindBySize(width int, height int) ([]model.Painting, error) {
	return g.paintings.FindBySize(width, height)
}

func (g *ArtGalleryApp) FindByGenre(genre string) ([]model.Painting, error) {
	return g.paintings.FindByGenre(genre)
}

func (g *ArtGalleryApp) FindByUserName(name string) ([]model.Painting, error) {
	return g.paintings.FindByUserName(name)
}

func (g *ArtGalleryApp) FindByUserEmail(email string) (*model.User, error) {
	return g.paintings.FindByUserEmail(email)
}

func (g *ArtGalleryApp) AddUser(user *model.User) error {
	return g.paintings.AddUser(user)
}

func (g *ArtGalleryApp) GetAllGenres() ([]string, error) {
	return g.paintings.GetAllGenres()
}

func (g *ArtGalleryApp) FindGenre(genre string) (*model.Genre, error) {
	return g.paintings.FindGenre(genre)
}

func (g *ArtGalleryApp) AddGenre(genre *model.Genre) error {
	return g.paintings.AddGenre(genre)
}

func (g *ArtGalleryApp) DeleteAllGenres() error {
	return g.paintings.DeleteAllGenres()
}

func (g *ArtGalleryApp) FindUserByID(id string) (*model.User, error) {
	return g.paintings.FindUserByID(id)
}

func (g *ArtGalleryApp) FindGenreByID(id string) (*model.Genre, error) {
	return g.paintings.FindGenreByID(id)
}

func (g *ArtGalleryApp) FindPaintingByID(id string) (*model.Painting, error) {
	return g.paintings.FindPaintingByID(id)
}

func (g *ArtGalleryApp) DeletePaintingByID(id string) error {
	return g.paintings.DeletePaintingByID(id)
}

func (g *ArtGalleryApp) FindGenreByName(name string) (*model.Genre, error) {
	return g.paintings.FindGenreByName(name)
}

func (g *ArtGalleryApp) FindUsersByFirstName(name string) ([]model.User, error) {
	return g.paintings.FindUsersByFirstName(name)
}

func (g *ArtGalleryApp) FindByUserNameAndPaintingTitle(name string, title string) ([]model.Painting, error) {
	return g.paintings.FindByUserNameAndPaintingTitle(name, title)
}

func (g *ArtGalleryApp) FindByUserNameAndGenre(name string, genreName string) ([]model.Painting, error) {
	return g.paintings.FindByUserNameAndGenre(name, genreName)
}

func (g *ArtGalleryApp) FindByPaintingTitleAndGenre(title string, genreName string) ([]model.Painting, error) {
	return g.paintings.FindByPaintingTitleAndGenre(title, genreName)
}

func (g *ArtGalleryApp) FindByUserNameAndPaintingTitleAndCenre(name string, title string, genreName string) ([]model.Painting, error) {
	return g.paintings.FindByUserNameAndPaintingTitleAndCenre(name, title, genreName)
}
