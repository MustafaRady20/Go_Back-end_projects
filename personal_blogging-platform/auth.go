package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIServer) Register(w http.ResponseWriter, r *http.Request) error {
	Newuser := &User{}
	err := json.NewDecoder(r.Body).Decode(&Newuser)
	if err != nil {
		return err
	}

	user := NewUser(Newuser.FirstName, Newuser.LastName, Newuser.Email, Newuser.Password)

	fmt.Println(user)
	if user.FirstName == "" || user.LastName == "" {
		return WriteJSON(w, http.StatusBadRequest, "User name is required")
	}

	emailValid := isEmailValid(user.Email)
	if !emailValid {
		return WriteJSON(w, http.StatusBadRequest, "Email is Not valid.")
	}
	passwordValid := isPasswordValid(user.Password)
	if !passwordValid {
		return WriteJSON(w, http.StatusBadRequest, "Password is Not valid.")
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, "can not generate password hash")
	}

	storedData := NewUser(user.FirstName, user.LastName, user.Email, string(hasedPassword))
	err = s.store.CreateUser(storedData)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}
	return nil

}
func (s *APIServer) Login(w http.ResponseWriter, r *http.Request) error {
	loginReq := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	user, err := s.store.GetUserByEmail(loginReq.Email)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}
	token, err := createJwt(user.ID)
	if err != nil {
		fmt.Println("Error creating JWT:", err)
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
func createJwt(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    id,
		"expiresAt": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET_KEY environment variable not set")
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func isAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := verifyToken(tokenString)
		if err != nil {
			log.Fatal("Not able to validate the token")
			WriteJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		if !token.Valid {
			log.Fatal("Not A Valid token")
			return
		}

		cliams := token.Claims.(jwt.MapClaims)

		id := cliams["userId"].(string)

		fmt.Println(id)
		// log.Fatal(id)

		ctx := context.WithValue(r.Context(), "userId", id)
		handler.ServeHTTP(w, r.WithContext(ctx))
		// handler(w, r)
	}

}

func verifyToken(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}
