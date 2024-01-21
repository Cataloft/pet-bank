package refresh

import "net/http"

type Request struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	AccessToken  string
	RefreshToken string
}

func Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
