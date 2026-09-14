package http

import (
	"errors"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

func writeGalleryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrTooManyImages):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrGalleryPostNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, media.ErrEmptyFile),
		errors.Is(err, media.ErrUnsupportedType),
		errors.Is(err, media.ErrFileTooLarge):
		response.Error(w, http.StatusBadRequest, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
