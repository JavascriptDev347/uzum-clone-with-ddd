package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

func writeCategoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyCategoryID),
		errors.Is(err, domain.ErrorEmptyCategoryName):
		response.Error(w, http.StatusBadRequest, err.Error()) // 400

	case errors.Is(err, domain.ErrCategoryNotFound):
		response.Error(w, http.StatusNotFound, err.Error()) // 404

	default:
		log.Println(err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
