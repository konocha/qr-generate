package store

import (
	//"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/konocha/qr-generate/internal/app/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	err := r.store.db.QueryRow("INSERT INTO users (email, encrypted_password) values (?, ?)", u.Email, u.EncryptedPassword)
	if err != nil {
		return err.Err()
	}
	r.store.db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error){
	u := &model.User{}

	err := r.store.db.QueryRow("SELECT * FROM users WHERE email=?", email).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	)
	if err != nil{
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Delete(email string) error{
	err := r.store.db.QueryRow("DELETE FROM users where email = ?", email)
	if err != nil{
		return err.Err()
	}

	return nil
}
