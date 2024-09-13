package store_test

import (
	"database/sql"
	"fmt"
	"testing"

	//"github.com/konocha/qr-generate/internal/app/model"
	"github.com/konocha/qr-generate/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func Test_Delete(t *testing.T) {
	db, err := store.TestDB()
	if err != nil {
		t.Fatal(err)
	}

	s := store.New(db)

	email1 := "Alex@gmail.com"
	email2 := "Pavel@gmail.com"

	err = s.User().Delete(email1)
	fmt.Println(err)
	assert.NoError(t, err)
	_, err = s.User().FindByEmail(email1)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	
	err = s.User().Delete(email2)
	assert.NoError(t, err)
	_, err = s.User().FindByEmail(email2)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
}
