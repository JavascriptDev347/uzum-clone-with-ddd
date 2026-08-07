package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

func writeIdentityError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrEmailAlreadyExists):
		response.Error(w, http.StatusConflict, err.Error()) // 409
	case errors.Is(err, application.ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, err.Error()) // 401
	case errors.Is(err, domain.ErrInvalidEmail):
		response.Error(w, http.StatusBadRequest, err.Error()) // 400
	default:
		log.Println(err)
		response.Error(w, http.StatusInternalServerError, "internal server error") // 500 - asl xatoni logga yozing, clientga chiqarmang!
	}
}
