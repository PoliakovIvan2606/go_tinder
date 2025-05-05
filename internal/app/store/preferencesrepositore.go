package store

import (
	"database/sql"
	"errors"
	"fmt"
	"tinder/internal/app/models"
)


type PreferencesRepository struct {
	store *Store
}

func (r *PreferencesRepository) Create(p *models.Preferences) (int, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}

	err := r.store.db.QueryRow(
		`INSERT INTO preferences (
			user_id,
			gender,
			age_from,
			age_to,
			radius
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			gender = $2,
			age_from = $3,
			age_to = $4,
			radius = $5
		RETURNING user_id`,
		p.User_id,
		p.Gender,
		p.Age_from,
		p.Age_to,
		p.Radius,
	).Scan(&p.User_id)
	return p.User_id, err
}

func (r *PreferencesRepository) GetByID(userID int) (*models.Preferences, error) {
	p := &models.Preferences{}
	
	err := r.store.db.QueryRow(
		`SELECT 
			user_id,
			gender,
			age_from,
			age_to,
			radius
		FROM preferences
		WHERE user_id = $1`,
		userID,
	).Scan(
		&p.User_id,
		&p.Gender,
		&p.Age_from,
		&p.Age_to,
		&p.Radius,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("preferences not found")
		}
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return p, nil
}
