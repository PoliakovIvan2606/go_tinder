package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
        "INSERT INTO users (name, email, password_hash, age, description, city, coordinates) VALUES ($1, $2, $3, $4, $5, $6, $7::geography) RETURNING id",
        u.Name,
        u.Email,
        u.PasswordHash,
        u.Age,
        u.Description,
        u.City,
        u.ToWKT(),
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
            ST_X(u.coordinates::geometry) AS longitude,
            ST_Y(u.coordinates::geometry) AS latitude,
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
        &u.Longitude,
        &u.Latitude,
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

func (r *UserRepository)IdFromUsers() (*sql.Rows, error){
    rows, err := r.store.db.Query("SELECT id FROM users")
    if err != nil {
        return nil, err
    }
    return rows, nil
}

func (r *UserRepository) IdPreferencesUser(userId int) (*sql.Rows, error) {
	// Получаем координаты и предпочтения пользователя
	var (
		lat, lon   float64
		ageFrom    int
		ageTo      int
		radius     int
	)

	err := r.store.db.QueryRow(`
		SELECT 
			ST_Y(coordinates::geometry), 
			ST_X(coordinates::geometry), 
			p.age_from, 
			p.age_to, 
			p.radius
		FROM users u
		JOIN preferences p ON p.user_id = u.id
		WHERE u.id = $1
	`, userId).Scan(&lat, &lon, &ageFrom, &ageTo, &radius)

	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}
	// Запрос пользователей в радиусе и нужном возрастном диапазоне
	rows, err := r.store.db.Query(`
		SELECT id
		FROM users
		WHERE
			id != $1 AND
			ST_DWithin(
				coordinates,
				ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography,
				$4
			)
			AND age BETWEEN $5 AND $6
		ORDER BY id
	`, userId, lon, lat, radius, ageFrom, ageTo)

	if err != nil {
		return nil, fmt.Errorf("query nearby users failed: %w", err)
	}

	return rows, nil
}