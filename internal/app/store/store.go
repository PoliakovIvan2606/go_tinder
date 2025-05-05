package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	db *sql.DB
	UserRepository *UserRepository
	PhotoRepository *PhotoRepository
	PreferencesRepository *PreferencesRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) User() *UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}

func (s *Store) Photo() *PhotoRepository {
	if s.PhotoRepository != nil {
		return s.PhotoRepository
	}

	s.PhotoRepository = &PhotoRepository{
		store: s,
	}

	return s.PhotoRepository
}

func (s *Store) Preferences() *PreferencesRepository {
	if s.PreferencesRepository != nil {
		return s.PreferencesRepository
	}

	s.PreferencesRepository = &PreferencesRepository{
		store: s,
	}

	return s.PreferencesRepository
}