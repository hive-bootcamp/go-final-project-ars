package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Arsadidas/go-final-project/pkg/helpers"
	"github.com/golang-jwt/jwt/v5"
)

type UserPassword struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	password := UserPassword{}
	envPassword := os.Getenv("TODO_PASSWORD")
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		helpers.WriteJSONError(w, "cannot get password")
		return
	}

	if password.Password != envPassword {
		helpers.WriteJSONError(w, "the password is incorrect")
		return
	}
	hashOfBytes := sha256.Sum256([]byte(password.Password))
	passwordHash := hex.EncodeToString(hashOfBytes[:])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passwordHash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		helpers.WriteJSONError(w, "cannot generate token")
		return
	}
	helpers.WriteJson(w, map[string]string{"token": signedToken})
}
