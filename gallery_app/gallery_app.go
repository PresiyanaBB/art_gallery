package gallery_app

import "art_gallery/model"

type PaintingRepository interface {
	GetAll() ([]model.Painting, error)
	Insert(painting *model.Painting) error
	DeleteAll() error
	FindByTitle(title string) ([]model.Painting, error)
	FindBySize(width int, height int) ([]model.Painting, error)
	FindByGenre(genre string) ([]model.Painting, error)
	FindByAuthor(name string) ([]model.Painting, error)
	SellPainting(id string) error
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

func (g *ArtGalleryApp) Add(p *model.Painting) error {
	return g.paintings.Insert(p)
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
