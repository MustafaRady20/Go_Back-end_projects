package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, "Cannot generate password hash")
	}

	storedData := NewUser(user.FirstName, user.LastName, user.Email, string(hashedPassword))
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
		"expiresAt": time.Now().Add(time.Minute * 15).Unix(),
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

		if !strings.HasPrefix(tokenString, "Bearer ") {
			WriteJSON(w, http.StatusUnauthorized, "Authorization token format must be Bearer <token>")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := verifyToken(tokenString)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				WriteJSON(w, http.StatusUnauthorized, "Token has expired")
				return
			}
			WriteJSON(w, http.StatusBadRequest, "Invalid token: "+err.Error())
			return
		}

		if !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		exp := int64(claims["expiresAt"].(float64))
		id := claims["userId"].(string)
		if time.Now().Unix() > exp {
			WriteJSON(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Log user ID
		fmt.Println("Authenticated user ID:", id)

		// Pass user ID in context to the handler
		ctx := context.WithValue(r.Context(), "userId", id)
		handler.ServeHTTP(w, r.WithContext(ctx))
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
