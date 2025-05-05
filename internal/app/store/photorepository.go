package store

import (
	"fmt"
	"strings"
	"tinder/internal/app/models"
)


type PhotoRepository struct {
	store *Store
}

func (r *PhotoRepository) Create(p *models.Photo) error {
	values := []interface{}{}
	placeholders := []string{}
	for i, url := range p.Photos {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, url)
	}
	query := fmt.Sprintf("INSERT INTO photos (user_id, photo_url) VALUES %s", strings.Join(placeholders, ","))
	args := append([]interface{}{p.User_id}, values...)
	_, err := r.store.db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}