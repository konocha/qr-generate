package apiserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/konocha/qr-generate/internal/app/auth"
	"github.com/konocha/qr-generate/internal/app/model"
	"github.com/konocha/qr-generate/internal/app/validation"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func (s *server) handleUserCreate() http.HandlerFunc {

	type UserCreate struct {
		Email    string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err, "error")
			return
		}
		defer r.Body.Close()

		var user UserCreate
		err = json.Unmarshal(body, &user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}

		err = validation.ValidateEmail(user.Email)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err, "error")
			return
		}
		err = validation.ValidatePassword(user.Password)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err, "error")
			return
		}

		u := &model.User{
			Email:             user.Email,
			EncryptedPassword: user.Password,
		}

		err = createFolder("./static/uploads/" + user.Email)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "ошибка создания папки")
			return
		}

		err = s.store.User().Create(u)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				s.error(w, r, http.StatusUnprocessableEntity, err, "Пользователь с такой почтой уже существует")
				return
			}

			s.error(w, r, http.StatusUnprocessableEntity, err, "error")
			return
		}

		json.NewEncoder(w).Encode("Пользователь с почтой " + user.Email + " успешно создан")
	}
}

func (s *server) handleUserAuth() http.HandlerFunc {
	type UserAuth struct {
		Email    string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err, "error")
			return
		}
		defer r.Body.Close()

		var user UserAuth

		err = json.Unmarshal(body, &user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}

		if user.Email == "" || user.Password == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("обязательно ввести почту и пароль"), "обязательно ввести почту и пароль")
			return
		}

		u, err := s.store.User().FindByEmail(user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusUnauthorized, err, "вы не зарегистрированы")
				return
			}

			s.error(w, r, http.StatusBadRequest, err, "error")
			return
		}

		if user.Email != u.Email || user.Password != u.EncryptedPassword {
			s.error(w, r, http.StatusUnauthorized, errors.New("неверное имя пользователя или пароль"), "неверное имя пользователя или пароль")
			return
		}

		token := auth.GenerateJWT(user.Email)
		json.NewEncoder(w).Encode(token)
	}
}

func (s *server) handleUserMe() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		if tokenString == ""{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("необходим токен авторизации"), "необходим токен авторизации")
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, err, "error")
			return
		}

		u, err := s.store.User().FindByEmail(claims.Email)
		if err != nil{
			if err == sql.ErrNoRows{
				s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("вы не зарегистрированы"), "вы не зарегистрированы")
				return
			}

			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(u.Email)
	}
}

func (s *server) handleUserDelete() http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request){
		tokenString := r.Header.Get("authorization")
		if tokenString == ""{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("необходим токен авторизации"), "необходим токен авторизации")
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, err, "Невалидный токен авторизации")
			return
		}

		_, err = s.store.User().FindByEmail(claims.Email)
		if err != nil{
			if err == sql.ErrNoRows{
				s.error(w, r, http.StatusBadRequest, err, "не существует пользователя с данной почтой")
				return
			}

			s.error(w, r, http.StatusInternalServerError, err, "error")
		}

		err = deleteFolder("./static/uploads/" + claims.Email)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "ошибка удаления папки")
			return
		}

		err = s.store.User().Delete(claims.Email)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}

		
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Пользователь с почтой " + claims.Email + " удален"})
	}
}

func (s *server) handleQRCreate() http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request){
		r.Header.Set("content-type", "application/json")
		tokenString := r.Header.Get("authorization")
		if tokenString == ""{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("необходим токен авторизации"), "необходим токен авторизации")
			return 
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, err, "невалидный токен авторизации")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil{
			s.error(w, r, http.StatusBadRequest, err, "error")
			return
		}
		defer r.Body.Close()

		var Qr model.QRcode
		err = json.Unmarshal(body, &Qr)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}


		qrCode, err := qr.Encode(Qr.Value, qr.M, qr.Auto)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}
		qrCode, err = barcode.Scale(qrCode, Qr.Width, Qr.Height)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}
		
		filename := "./static/uploads/" + claims.Email + "/" + Qr.Value + ".png"
		file, err := os.Create(filename)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "ошибка создания файла")
		}
		defer file.Close()

		png.Encode(file, qrCode)

		file, err = os.Open(filename)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "ошибка чтения файла")
			return
		}
		data, err := io.ReadAll(file)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "ошибка чтения файла")
		}
		
		s.store.QrCode().Create(claims.Email, Qr.Value)
		w.Header().Set("content-type", "image/png")
		w.Write(data)
	}
}

func (s *server) handleQRDelete() http.HandlerFunc{
	type qrValue struct{
		Value string
	}
	
	return func(w http.ResponseWriter, r *http.Request){
		r.Header.Set("content-type", "application/json")
		tokenString := r.Header.Get("authorization")
		if tokenString == ""{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("необходим токен авторизации"), "необходим токен авторизации")
			return 
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, err, "невалидный токен авторизации")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
		}
		defer r.Body.Close()

		var qrcode qrValue
		err = json.Unmarshal(body, &qrcode)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
		}

		filename := "./static/uploads/" + claims.Email + "/" + qrcode.Value + ".png"

		err = os.Remove(filename)
		if err != nil{
			s.error(w, r, http.StatusBadRequest, err, "не существует")
			return
		}

		json.NewEncoder(w).Encode("qr-код со значением " + qrcode.Value + " удален")
		s.store.QrCode().Delete(claims.Email, qrcode.Value)
	}
}

func (s * server) handleQRGet() http.HandlerFunc{
	type qrValue struct{
		Value string
	}

	return func(w http.ResponseWriter, r *http.Request){
		r.Header.Set("content-type", "application-json")

		r.Header.Set("content-type", "application/json")
		tokenString := r.Header.Get("authorization")
		if tokenString == ""{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("необходим токен авторизации"), "необходим токен авторизации")
			return 
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, err, "невалидный токен авторизации")
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
		}
		defer r.Body.Close()

		var qrcode qrValue
		err = json.Unmarshal(body, &qrcode)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
		}

		filename := "./static/uploads/" + claims.Email + "/" + qrcode.Value + ".png"
		file, err := os.Open(filename)
		if err != nil{
			s.error(w, r, http.StatusBadRequest, err, "не существует")
			return
		}
		data, err := io.ReadAll(file)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}

		w.Header().Set("content-type", "image/png")
		w.Write(data)
	}
}

func (s *server) handleQRAll() http.HandlerFunc{
	type QrResponse struct{
		Qr_values []string
	}

	return func(w http.ResponseWriter, r *http.Request){
		r.Header.Set("content-type", "application-json")

		r.Header.Set("content-type", "application/json")
		tokenString := r.Header.Get("authorization")
		if tokenString == ""{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, errors.New("необходим токен авторизации"), "необходим токен авторизации")
			return 
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil{
			s.error(w, r, http.StatusNetworkAuthenticationRequired, err, "невалидный токен авторизации")
			return
		}

		values, err := s.store.QrCode().GetAll(claims.Email)
		if err != nil{
			s.error(w, r, http.StatusInternalServerError, err, "error")
			return
		}

		Qr := QrResponse{
			Qr_values: values,
		}

		json.NewEncoder(w).Encode(Qr)
	}
}


func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error(), "message": msg})
}

func createFolder(path string) error{
	return os.Mkdir(path, 0755)
}

func deleteFolder(path string) error{
	return os.RemoveAll(path)
}
