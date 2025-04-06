package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gallery struct {
	ID        int
	UserID    int
	Title     string
	Published bool
}

type GalleryService struct {
	DB *sql.DB
}

func (gs *GalleryService) Create(title string, userID int, published bool) (*Gallery, error) {
	gallery := Gallery{
		Title:     title,
		UserID:    userID,
		Published: published,
	}
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (user_id, title, published)
		VALUES ($2, $1, $3) RETURNING id;`, title, userID, published)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}
	row := gs.DB.QueryRow(`
		SELECT title, user_id, published
		FROM galleries WHERE id = $1;`, id)
	err := row.Scan(&gallery.Title, &gallery.UserID, &gallery.Published)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, title, published
		FROM galleries
		WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err := rows.Scan(&gallery.ID, &gallery.Title, &gallery.Published)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	res, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $1, published = $2
		WHERE id = $3;`, gallery.Title, gallery.Published, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	fmt.Println(res)
	return nil
}

func (gs *GalleryService) Delete(id int) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}
