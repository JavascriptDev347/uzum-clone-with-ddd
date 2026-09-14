package http

import (
	"errors"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

var ErrUnauthorized = errors.New("review: unauthorized")

func writeReviewError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRating):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotEligibleForReview):
		response.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrReviewAlreadyExists):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrReviewNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrUnauthorized):
		response.Error(w, http.StatusUnauthorized, "unauthorized")
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
