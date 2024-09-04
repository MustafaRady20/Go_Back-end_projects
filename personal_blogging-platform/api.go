package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ServerPort string
	store      Storage
}

func newAPIServer(serverPort string, store Storage) *APIServer {
	return &APIServer{
		ServerPort: serverPort,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/register", makeHTTPHandleFunc(s.Register)).Methods("POST")
	router.HandleFunc("/login", makeHTTPHandleFunc(s.Login)).Methods("POST")
	router.HandleFunc("/user", isAuthorized(makeHTTPHandleFunc(s.getUser))).Methods("GET")
	fmt.Printf("Server is running on port %v", s.ServerPort)
	http.ListenAndServe(s.ServerPort, router)

}

func (s *APIServer) getUser(w http.ResponseWriter, r *http.Request) error {
	// r.Context().Value("userId")
	fmt.Println(r.Context().Value("userId").(string))
	return nil
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}
