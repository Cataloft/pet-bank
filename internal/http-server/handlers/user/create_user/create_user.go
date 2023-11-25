package create_user

import (
	"bank-api/internal/model"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type Response struct {
	Status int
	Error  string
}

//type Service struct {
//	userExistence UserExistence
//	userSaver     UserSaver
//}
//
//type UserExistence interface {
//	UserExists(email string) (bool, error)
//}

type UserSaver interface {
	SaveUser(user *model.User) error
}

// user UserSaver

func CreateUser(userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		var resp Response

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			resp.Status = http.StatusBadRequest
			resp.Error = "failed to decode request body"
			log.Println(resp, err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
			FullName: req.FullName,
		}

		if err := userSaver.SaveUser(u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to save user:", err)
			return
		}

		u.Sanitize()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(u.Email))

		log.Println(resp.Status, "CREATED")

	}
}
