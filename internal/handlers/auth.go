package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Enotisi/go_final_project/internal/config"
	"github.com/Enotisi/go_final_project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const sercretKey = "MySecretKey"

func SignHandle(w http.ResponseWriter, r *http.Request) {

	passwordData := models.PasswordRequest{}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &passwordData); err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	if passwordData.Password != config.Conf.Password {
		http.Error(w, createJsonResponse("error", "неверный пароль"), http.StatusUnauthorized)
		return
	}

	tokeLifetime := time.Now().Add(time.Duration(config.Conf.TokenLifeTime) * time.Hour)
	passwordData.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(tokeLifetime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, passwordData)
	tokenStr, err := token.SignedString([]byte(sercretKey))

	if err != nil {
		http.Error(w, createJsonResponse("error", "неудалось создать токен"), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: tokeLifetime,
	})

	resp := createJsonResponse("token", tokenStr)
	if _, err := w.Write([]byte(resp)); err != nil {
		log.Printf("error return responce at GetTaskHandle: %s\n", err.Error())
	}
}

func WebHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path

	if (strings.Contains(url, ".html") || url == "/") && config.Conf.Password != "" {
		http.ServeFile(w, r, config.Conf.WebPath+"/login.html")
	} else {
		http.FileServer(http.Dir(config.Conf.WebPath)).ServeHTTP(w, r)
	}
}

func MiddlewareHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Conf.Password != "" {
			tokenStr, err := checkCookie(r)
			if err != nil {
				http.Error(w, createJsonResponse("error", err.Error()), http.StatusUnauthorized)
			}
			fmt.Println(tokenStr)

		}
		next.ServeHTTP(w, r)
	})
}

func checkCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			return "", errors.New("требуется аутентификация")
		}
		return "", err
	}

	return cookie.Value, nil
}
