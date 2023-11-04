package register_user

import (
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type Request struct {
	Email    string
	Password string
}

type Response struct {
	Status int
	Error  string
}

type UserSaver interface {
	SaveUser(string, string, []byte) error
}

func RegisterUser(saveUser UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var resp Response

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			resp.Status = http.StatusBadRequest
			resp.Error = "failed to decode request body"

			log.Println(resp, err)

			return
		}

		log.Println("request body decoded")

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("failed to hash password")
			return
		}

		err = saveUser.SaveUser(req.Email, req.Password, encryptedPassword)
		if err != nil {
			log.Println("failed to exec sql to db", err)
			return
		}
	}
}
