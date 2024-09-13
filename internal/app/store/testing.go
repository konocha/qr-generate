package store

import(
	"database/sql"
)

func TestDB() (*sql.DB, error){
	db, err := sql.Open("mysql", "username:password@tcp(0.0.0.0:3306)/nameOfDatabase")
	if err != nil{
		return nil, err
	}

	err = db.Ping()
	if err != nil{
		return nil, err
	}

	return db, nil
}