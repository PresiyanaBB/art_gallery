package gallery_app

import model "art_gallery/models"

type PaintingRepository interface {
	GetAll() ([]model.Painting, error)
	AddPainting(painting *model.Painting) error
	DeleteAll() error
	FindByTitle(title string) ([]model.Painting, error)
	FindBySize(width int, height int) ([]model.Painting, error)
	FindByGenre(genre string) ([]model.Painting, error)
	FindByAuthor(name string) ([]model.Painting, error)
	SellPainting(id string) error
	FindByUserEmail(email string) (*model.User, error)
	AddUser(u *model.User) error
	GetAllGenres() ([]string, error)
	FindGenre(genre string) (*model.Genre, error)
	AddGenre(genre *model.Genre) error
	RemoveAllGenres() error
	FindGenreByID(id string) (*model.Genre, error)
	FindUserByID(id string) (*model.User, error)
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

func (g *ArtGalleryApp) FindByAuthor(name string) ([]model.Painting, error) {
	return g.paintings.FindByAuthor(name)
}

func (g *ArtGalleryApp) SellPainting(id string) error {
	return g.paintings.SellPainting(id)
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

func (g *ArtGalleryApp) RemoveAllGenres() error {
	return g.paintings.RemoveAllGenres()
}

func (g *ArtGalleryApp) FindUserByID(id string) (*model.User, error) {
	return g.paintings.FindUserByID(id)
}

func (g *ArtGalleryApp) FindGenreByID(id string) (*model.Genre, error) {
	return g.paintings.FindGenreByID(id)
}
