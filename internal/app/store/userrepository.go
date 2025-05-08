package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"tinder/internal/app/models"
)


type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *models.UserCreate) (int, error) {
	if err := u.Validate(); err != nil {
		return 0, err
	}

    if err := u.HashPassword(); err != nil {
        return 0, err
    }
	
	var id int
    err := r.store.db.QueryRow(
        "INSERT INTO users (name, email, password_hash, age, description, city, coordinates) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
        u.Name,
        u.Email,
        u.PasswordHash,
        u.Age,
        u.Description,
        u.City,
        u.Coordinates,
    ).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}


func (r *UserRepository) UserById(id int) (*models.User, error) {
    u := &models.User{}
    var photosJSON []byte

    query := `
        SELECT
            u.id,
            u.name,
            u.email,
            u.password_hash,
            u.age,
            u.description,
            u.city,
            u.coordinates,
            u.created_at,
            COALESCE(
                json_agg(p.photo_url) FILTER (WHERE p.photo_url IS NOT NULL),
                '[]'
            ) AS photos
        FROM users u
        LEFT JOIN photos p ON p.user_id = u.id
        WHERE u.id = $1
        GROUP BY u.id;
    `

    err := r.store.db.QueryRow(query, id).Scan(
        &u.ID,
        &u.Name,
        &u.Email,
        &u.Password,
        &u.Age,
        &u.Description,
        &u.City,
        &u.Coordinates,
        &u.CreatedAt,
        &photosJSON,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // или своя логика, если пользователь не найден
        }
        return nil, err
    }

    // Парсим JSON массив в слайс строк
    if err := json.Unmarshal(photosJSON, &u.Photos); err != nil {
        return nil, err
    }

    return u, nil
}


func (r *UserRepository) IdAndPaswordByEmail(email string) (string, string, error) {
    var id string
    var passwordHash string
    err := r.store.db.QueryRow(`SELECT id, password_hash FROM users WHERE email = $1`, email).Scan(
        &id,
        &passwordHash,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", "", err
        }
        return "", "", err
    }

    return id, passwordHash, nil
}