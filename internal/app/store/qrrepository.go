package store

import "errors"

type QrRepository struct {
	store *Store
}

func (r *QrRepository) Create(email string, value string) error{
	var user_id int
	err := r.store.db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&user_id)
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec("INSERT INTO qr (qr_value, user_id) values (?, ?)", value, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (r * QrRepository) Delete(email string, value string) error{
	var user_id int
	err := r.store.db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&user_id)
	if err != nil {
		return err
	}

	row, err := r.store.db.Exec("DELETE FROM qr WHERE qr_value = ? AND user_id = ?", value, user_id)
	if err != nil{
		return err
	}
	if res, _ := row.RowsAffected(); res == 0{
		return errors.New("not exist")
	}

	return nil
} 

func (r *QrRepository) GetAll(email string) ([]string, error){
	var values []string
	var user_id int
	err := r.store.db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&user_id)
	if err != nil {
		return []string{}, err
	}

	rows, err := r.store.db.Query("SELECT qr_value FROM qr WHERE user_id = ?", user_id)
	if err != nil{
		return []string{}, err
	}
	
	for rows.Next(){
		var value string
		rows.Scan(&value)
		values = append(values, value)
	}

	return values, nil
}
