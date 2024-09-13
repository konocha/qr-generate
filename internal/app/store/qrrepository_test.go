package store_test

import (
	"testing"
	"errors"

	"github.com/konocha/qr-generate/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T){
	db, err := store.TestDB()
	if err != nil{
		t.Fatal(err)
	}

	s := store.New(db)

	//Not existing email
	emailNotExist := "someEmail@gmail.com"
	err = s.QrCode().Create(emailNotExist, "someValue")
	assert.Error(t, err)

	//Existing email
	emailExist := "hello@gmail.com"
	err = s.QrCode().Create(emailExist, "someValue")
	assert.NoError(t, err)
} 

func Test_DeleteQR(t *testing.T){
	db, err := store.TestDB()
	if err != nil{
		t.Fatal(err)
	}

	s := store.New(db)

	//Not existing email and value
	emailNotExist := "someEmail@gmail.com"
	valueNotExist := "someValue"
	err = s.QrCode().Delete(emailNotExist, valueNotExist)
	assert.Error(t, err)

	//Not existing email and existing value
	emailNotExist = "someEmail@gmail.com"
	valueExist := "postman"
	err  = s.QrCode().Delete(emailNotExist, valueExist)
	assert.Error(t, err)

	//Existing email and not existing value
	emailExist := "hello@gmail.com"
	valueNotExist = "someValue123"
	err = s.QrCode().Delete(emailExist, valueNotExist)
	if err.Error() != errors.New("not exist").Error(){
		t.Error(err)
	}

	//Existing email and existing value
	emailExist = "hello@gmail.com"
	valueExist = "someValue" 
	err = s.QrCode().Delete(emailExist, valueExist)
	assert.NoError(t, err)

}