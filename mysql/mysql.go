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
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
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
	_, err := r.client.Exec("INSERT INTO paintings(id, title, description, mime_type, data , author, date_of_publication, width, height, genre,price) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		p.ID, p.Title, p.Description, p.MIMEType, p.Data, p.Author.ID, p.DateOfPublication.String(), p.Width, p.Height, p.Genre.ID, p.Price)

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
		var date string
		var authID string
		var genreID string
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
		u, _ := r.FindUserByID(authID)
		g, _ := r.FindGenreByID(genreID)
		result.Author = *u
		result.Genre = *g
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
		var date string
		var authID string
		var genreID string
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
		u, _ := r.FindUserByID(authID)
		g, _ := r.FindGenreByID(genreID)
		result.Author = *u
		result.Genre = *g
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

	g, _ := r.FindGenre(genre)
	rows, err := r.client.Query("SELECT * FROM paintings WHERE genre = ?", g.ID)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		var date string
		var authID string
		var genreID string
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
		u, _ := r.FindUserByID(authID)
		g, _ := r.FindGenreByID(genreID)
		result.Author = *u
		result.Genre = *g
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) FindByUserName(name string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	users, _ := r.FindUsersByFirstName(name)

	for _, v := range users {
		query := fmt.Sprintf("SELECT * FROM paintings WHERE author like '%v%v%v'", "%", v.ID, "%")
		rows, err := r.client.Query(query)
		if err != nil {
			return nil, fmt.Errorf("mysql query failure: %w", err)
		}
		for rows.Next() {
			var result model.Painting
			var date string
			var authID string
			var genreID string
			rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
			result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
			u, _ := r.FindUserByID(authID)
			g, _ := r.FindGenreByID(genreID)
			result.Author = *u
			result.Genre = *g
			paintings = append(paintings, result)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating paintings: %w", err)
		}
		rows.Close()
	}

	return paintings, nil
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
		result.DateOfRegistration, _ = time.Parse("2006-01-02 15:04:05.9999999", dateOfReg)
		users = &result
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
		users.DateOfRegistration, _ = time.Parse("2006-01-02 15:04:05.9999999", dateOfReg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

func (r *MySQLRepository) FindPaintingByID(id string) (*model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var painting model.Painting

	rows, err := r.client.Query("SELECT * FROM paintings WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		var date string
		var authID string
		var genreID string
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
		u, _ := r.FindUserByID(authID)
		g, _ := r.FindGenreByID(genreID)
		result.Author = *u
		result.Genre = *g
		painting = result
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return &painting, nil
}

func (r *MySQLRepository) DeletePaintingByID(id string) error {
	if r.client == nil {
		return fmt.Errorf("mysql repository is not initilized")
	}
	_, err := r.client.Exec("DELETE FROM paintings WHERE id = ?", id)

	return err
}

func (r *MySQLRepository) FindGenreByName(name string) (*model.Genre, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}

	rows, err := r.client.Query("SELECT * FROM genres WHERE name = ?", name)
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

func (r *MySQLRepository) FindUsersByFirstName(name string) ([]model.User, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE first_name like '%v%v%v'", "%", name, "%")
	rows, err := r.client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var result model.User
		var dateOfReg string
		rows.Scan(&result.ID, &result.FirstName, &result.LastName, &dateOfReg, &result.Email, &result.Password)
		result.DateOfRegistration, _ = time.Parse("2006-01-02 15:04:05.9999999", dateOfReg)
		users = append(users, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}
	return users, nil
}

func (r *MySQLRepository) FindByUserNameAndPaintingTitle(name string, title string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	users, _ := r.FindUsersByFirstName(name)

	for _, v := range users {
		query := fmt.Sprintf("SELECT * FROM paintings WHERE author like '%v%v%v' and title like '%v%v%v'", "%", v.ID, "%", "%", title, "%")
		rows, err := r.client.Query(query)
		if err != nil {
			return nil, fmt.Errorf("mysql query failure: %w", err)
		}
		for rows.Next() {
			var result model.Painting
			var date string
			var authID string
			var genreID string
			rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
			result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
			u, _ := r.FindUserByID(authID)
			g, _ := r.FindGenreByID(genreID)
			result.Author = *u
			result.Genre = *g
			paintings = append(paintings, result)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating paintings: %w", err)
		}
		rows.Close()
	}

	return paintings, nil
}

func (r *MySQLRepository) FindByUserNameAndGenre(name string, genreName string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	genre, _ := r.FindGenreByName(genreName)
	users, _ := r.FindUsersByFirstName(name)

	for _, v := range users {
		query := fmt.Sprintf("SELECT * FROM paintings WHERE author like '%v%v%v' and genre like '%v%v%v'", "%", v.ID, "%", "%", genre.ID, "%")
		rows, err := r.client.Query(query)
		if err != nil {
			return nil, fmt.Errorf("mysql query failure: %w", err)
		}
		for rows.Next() {
			var result model.Painting
			var date string
			var authID string
			var genreID string
			rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
			result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
			u, _ := r.FindUserByID(authID)
			g, _ := r.FindGenreByID(genreID)
			result.Author = *u
			result.Genre = *g
			paintings = append(paintings, result)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating paintings: %w", err)
		}
		rows.Close()
	}

	return paintings, nil
}

func (r *MySQLRepository) FindByPaintingTitleAndGenre(title string, genreName string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	genre, _ := r.FindGenreByName(genreName)

	query := fmt.Sprintf("SELECT * FROM paintings WHERE title like '%v%v%v' and genre like '%v%v%v'", "%", title, "%", "%", genre.ID, "%")
	rows, err := r.client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("mysql query failure: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result model.Painting
		var date string
		var authID string
		var genreID string
		rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
		result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
		u, _ := r.FindUserByID(authID)
		g, _ := r.FindGenreByID(genreID)
		result.Author = *u
		result.Genre = *g
		paintings = append(paintings, result)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating paintings: %w", err)
	}
	return paintings, nil
}

func (r *MySQLRepository) FindByUserNameAndPaintingTitleAndCenre(name string, title string, genreName string) ([]model.Painting, error) {
	if r.client == nil {
		return nil, fmt.Errorf("mysql repository is not initilized")
	}
	var paintings []model.Painting

	genre, _ := r.FindGenreByName(genreName)
	users, _ := r.FindUsersByFirstName(name)

	for _, v := range users {
		query := fmt.Sprintf("SELECT * FROM paintings WHERE author like '%v%v%v' and genre like '%v%v%v' and title like '%v%v%v'", "%", v.ID, "%", "%", genre.ID, "%", "%", title, "%")
		rows, err := r.client.Query(query)
		if err != nil {
			return nil, fmt.Errorf("mysql query failure: %w", err)
		}
		for rows.Next() {
			var result model.Painting
			var date string
			var authID string
			var genreID string
			rows.Scan(&result.ID, &result.Title, &result.Description, &result.MIMEType, &result.Data, &authID, &date, &result.Width, &result.Height, &genreID, &result.Price)
			result.DateOfPublication, _ = time.Parse("2006-01-02 15:04:05.9999999", date)
			u, _ := r.FindUserByID(authID)
			g, _ := r.FindGenreByID(genreID)
			result.Author = *u
			result.Genre = *g
			paintings = append(paintings, result)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating paintings: %w", err)
		}
		rows.Close()
	}

	return paintings, nil
}
