package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

func writeProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyProductID),
		errors.Is(err, domain.ErrEmptyProductName),
		errors.Is(err, domain.ErrEmptyCategoryID),
		errors.Is(err, domain.ErrNegativeAmount),
		errors.Is(err, domain.ErrInvalidCurrency),
		errors.Is(err, domain.ErrCurrencyMismatch):
		response.Error(w, http.StatusBadRequest, err.Error()) // 400

	default:
		log.Println(err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
