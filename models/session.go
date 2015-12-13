package models

import (
	"net/http"
)

func MySession(r *http.Request) (value string) {
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return value
	}
	value = cookie.Value
	return value
}
