package main

import (
	"encoding/json"
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
	router.HandleFunc("/blogs", isAuthorized(makeHTTPHandleFunc(s.getArticles))).Methods("GET")
	router.HandleFunc("/article", isAuthorized(makeHTTPHandleFunc(s.createNewArticle))).Methods("POST")
	router.HandleFunc("/UpdateArticle", isAuthorized(makeHTTPHandleFunc(s.updateArticle))).Methods("POST")


	fmt.Printf("Server is running on port %v", s.ServerPort)
	http.ListenAndServe(s.ServerPort, router)

}

func (s *APIServer) updateArticle(w http.ResponseWriter, r *http.Request) error {

	var updated UpdatedArticle
	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	userId := r.Context().Value("userId").(string)
	err = s.store.UpdateArticle(updated, userId)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	return nil
}

func (s *APIServer) getArticles(w http.ResponseWriter, r *http.Request) error {
	articles, err := s.store.getAllArticles()
	if err != nil {
		return err
	}
	WriteJSON(w, http.StatusOK, map[string][]*Article{
		"artilces": articles,
	})
	return nil
}

func (s *APIServer) createNewArticle(w http.ResponseWriter, r *http.Request) error {
	var article Article
	err := json.NewDecoder(r.Body).Decode(&article)
	fmt.Println(article)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	userId := r.Context().Value("userId").(string)
	fmt.Println(userId)
	if userId == "" {
		return WriteJSON(w, http.StatusBadRequest, "Log in again user Not found")
	}

	err = s.store.CreateArticle(userId, article)
	fmt.Println(err)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	fmt.Println("success")
	return nil
}
