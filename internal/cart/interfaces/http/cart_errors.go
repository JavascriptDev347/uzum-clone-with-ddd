package http

import (
	"errors"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/domain"
	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

var ErrUnauthorized = errors.New("cart: unauthorized")

func writeCartError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidQuantity):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrItemNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrEmptyProductID):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmptyUserID):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInsufficientStock):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, catalog.ErrProductNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrUnauthorized):
		response.Error(w, http.StatusUnauthorized, "unauthorized")
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
