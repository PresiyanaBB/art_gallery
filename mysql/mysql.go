package mysql

import (
	model "art_gallery/models"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLRepository struct {
	opts   MySQLOptions
	client *sql.DB
}

type MySQLOptions struct {
	URI string
}

func New(opts MySQLOptions) *MySQLRepository {
	return &MySQLRepository{
		opts:   opts,
		client: nil,
	}
}

func (r *MySQLRepository) Init() error {
	var err error
	r.client, err = sql.Open("mysql", r.opts.URI)
	return err
}

func (r *MySQLRepository) GetAll() ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	rows, err := r.client.Query("SELECT * FROM paintings")
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		var date string
		var authID string
		var genreID string
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.Src, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("Jan 2, 2006 at 3:04pm (MST)", date)
		u, _ := r.FindUserByID(authID)
		g, _ := r.FindGenreByID(genreID)
		result.Author = *u
		result.Genre = *g
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating images: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) AddPainting(p *model.Painting) error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}
	_, err := r.client.Exec("INSERT INTO paintings(id, title, description, src, author, date_of_publication, width, height, genre,price) VALUES (?,?,?,?,?,?,?,?,?,?)",
		p.ID, p.Title, p.Description, p.Src, p.Author.ID, p.DateOfPublication.String(), p.Width, p.Height, p.Genre.ID, p.Price)

	return err
}

func (r *MySQLRepository) DeleteAll() error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}
	_, err := r.client.Exec("TRUNCATE TABLE paintings")

	return err
}

func (r *MySQLRepository) FindByTitle(title string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	rows, err := r.client.Query("SELECT * FROM paintings WHERE title = ?", title)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.Src, &result.Author, &result.DateOfPublication, &result.Width, &result.Height, &result.Genre)
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) FindBySize(width int, height int) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	rows, err := r.client.Query("SELECT * FROM paintings WHERE width = ? and height = ?", width, height)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.Src, &result.Author, &result.DateOfPublication, &result.Width, &result.Height, &result.Genre)
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) FindByGenre(genre string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	rows, err := r.client.Query("SELECT * FROM paintings WHERE genre = ?", genre)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.Src, &result.Author, &result.DateOfPublication, &result.Width, &result.Height, &result.Genre)
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) FindByAuthor(name string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	rows, err := r.client.Query("SELECT * FROM paintings WHERE author like '%?%'", name)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.Src, &result.Author, &result.DateOfPublication, &result.Width, &result.Height, &result.Genre)
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) SellPainting(id string) error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}

	_, err := r.client.Exec("DELETE FROM paintings WHERE id = ?", id)

	return err
}

func (r *MySQLRepository) FindByUserEmail(email string) (*model.User, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}

	rows, err := r.client.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()

	var users *model.User

	for rows.Next() {
		var result model.User
		var dateOfReg string
		rows.Scan(&result.ID, &result.FirstName, &result.LastName, &dateOfReg, &result.Email, &result.Password)
		users = &result
		users.DateOfRegistration, _ = time.Parse("Jan 2, 2006 at 3:04pm (MST)", dateOfReg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

func (r *MySQLRepository) AddUser(u *model.User) error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}

	_, err := r.client.Exec("INSERT INTO users (id, first_name, last_name, date_of_registration, email, password) VALUES (?,?,?,?,?,?)",
		u.ID, u.FirstName, u.LastName, u.DateOfRegistration, u.Email, u.Password)

	return err
}

func (r *MySQLRepository) GetAllGenres() ([]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var genres []string

	rows, err := r.client.Query("SELECT name FROM genres")
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result string
		rows.Scan(&result)
		genres = append(genres, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return genres, nil
}

func (r *MySQLRepository) FindGenre(genre string) (*model.Genre, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var res model.Genre

	rows, err := r.client.Query("SELECT * FROM genres WHERE name = ?", genre)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return &res, nil
}

func (r *MySQLRepository) AddGenre(genre *model.Genre) error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}

	_, err := r.client.Exec("INSERT INTO genres (id, name) VALUES (?,?)",
		genre.ID, genre.Name)

	return err
}

func (r *MySQLRepository) RemoveAllGenres() error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}

	_, err := r.client.Exec("DELETE FROM genres")
	return err
}

func (r *MySQLRepository) FindGenreByID(id string) (*model.Genre, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}

	rows, err := r.client.Query("SELECT * FROM genres WHERE id = ?", id)
	var res model.Genre
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&res.ID, &res.Name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return &res, nil
}

func (r *MySQLRepository) FindUserByID(id string) (*model.User, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}

	rows, err := r.client.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()

	var users *model.User

	for rows.Next() {
		var result model.User
		var dateOfReg string
		rows.Scan(&result.ID, &result.FirstName, &result.LastName, &dateOfReg, &result.Email, &result.Password)
		users = &result
		users.DateOfRegistration, _ = time.Parse("Jan 2, 2006 at 3:04pm (MST)", dateOfReg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}
