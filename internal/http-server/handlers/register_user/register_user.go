package register_user

import (
	"bank-api/internal/model"
	"github.com/go-chi/render"
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
	SaveUser(user *model.User) error
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

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		log.Println("request body decoded")

		//req.EncryptedPassword, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		//if err != nil {
		//	log.Println("failed to hash password")
		//	return
		//}

		err = saveUser.SaveUser(u)
		if err != nil {
			log.Println("failed to exec sql to db", err)
			return
		}

		u.Sanitize()

	}
}
