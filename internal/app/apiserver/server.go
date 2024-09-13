package apiserver

import (

	"github.com/gorilla/mux"
	"github.com/konocha/qr-generate/internal/app/store"
)

type server struct{
	router *mux.Router
	store  *store.Store
}

func newServer(store *store.Store) *server{
	s := &server{
		router: mux.NewRouter(),
		store: store,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter(){
	s.router.HandleFunc("/user/create", s.handleUserCreate()).Methods("POST")
	s.router.HandleFunc("/user/auth", s.handleUserAuth()).Methods("POST") 
	s.router.HandleFunc("/user/me", s.handleUserMe()).Methods("GET")
	s.router.HandleFunc("/user/delete", s.handleUserDelete()).Methods("DELETE")
	s.router.HandleFunc("/qr/create", s.handleQRCreate()).Methods("POST")
	s.router.HandleFunc("/qr/delete", s.handleQRDelete()).Methods("DELETE")
	s.router.HandleFunc("/qr/get", s.handleQRGet()).Methods("GET")
	s.router.HandleFunc("/qr/all", s.handleQRAll()).Methods("GET")
}
