package store

import(
	"database/sql"
)

func TestDB() (*sql.DB, error){
	db, err := sql.Open("mysql", "root:Sofa=22082014@tcp(0.0.0.0:3306)/testDatabase")
	if err != nil{
		return nil, err
	}

	err = db.Ping()
	if err != nil{
		return nil, err
	}

	return db, nil
}