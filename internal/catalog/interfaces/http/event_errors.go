package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

func writeEventError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyEventID),
		errors.Is(err, domain.ErrEmptyEventTitle),
		errors.Is(err, domain.ErrEmptyCategoryID),
		errors.Is(err, domain.ErrCategoryNotFound), // category_id noto'g'ri/mavjud emas - so'rov xatosi
		errors.Is(err, media.ErrEmptyFile),
		errors.Is(err, media.ErrFileTooLarge),
		errors.Is(err, media.ErrUnsupportedType):
		response.Error(w, http.StatusBadRequest, err.Error()) // 400

	case errors.Is(err, domain.ErrEventNotFound):
		response.Error(w, http.StatusNotFound, err.Error()) // 404

	default:
		log.Println(err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
