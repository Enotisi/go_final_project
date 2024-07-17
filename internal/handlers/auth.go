package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

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

	claims := jwt.MapClaims{
		"hash": getHash(passwordData.Password),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(sercretKey))

	if err != nil {
		http.Error(w, createJsonResponse("error", "неудалось создать токен"), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenStr,
	})

	resp := createJsonResponse("token", tokenStr)
	if _, err := w.Write([]byte(resp)); err != nil {
		log.Printf("error return responce at GetTaskHandle: %s\n", err.Error())
	}
}

func WebHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir(config.Conf.WebPath)).ServeHTTP(w, r)
}

func MiddlewareHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Conf.Password != "" {
			tokenStr, err := checkCookie(r)
			if err != nil {
				http.Error(w, createJsonResponse("error", err.Error()), http.StatusUnauthorized)
				return
			}

			valid, err := checkToken(tokenStr)
			if err != nil {
				http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
				return
			}

			if !valid {
				http.Error(w, createJsonResponse("error", "неверный токен"), http.StatusUnauthorized)
				return
			}
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

func checkToken(tokenStr string) (bool, error) {
	jwtToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(sercretKey), nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	hash, ok := claims["hash"]
	if !ok {
		return false, nil
	}

	hashStr, ok := hash.(string)
	if !ok {
		return false, nil
	}

	if hashStr != getHash(config.Conf.Password) {
		return false, nil
	}

	return true, nil
}

func getHash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
