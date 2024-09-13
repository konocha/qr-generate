package store

import (
	"database/sql"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	qrRepository *QrRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() *UserRepository{
	if s.userRepository != nil{
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store : s,
	}

	return s.userRepository
}

func (s *Store) QrCode() *QrRepository{
	if s.qrRepository != nil{
		return s.qrRepository
	}

	s.qrRepository = &QrRepository{
		store : s,
	}

	return s.qrRepository
}
