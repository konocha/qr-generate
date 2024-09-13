package apiserver

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/konocha/qr-generate/internal/app/store"
)

func Start(config *Config) error{
	db, err := newDB(config)
	if err != nil {
		return err
	}
	defer db.Close()

	str := store.New(db)

	s := newServer(str)

	return http.ListenAndServe(config.BindAddress, s.router)
}

func newDB(config *Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DataBaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
