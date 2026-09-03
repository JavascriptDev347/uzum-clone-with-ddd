package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

func writeProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyProductID),
		errors.Is(err, domain.ErrEmptyProductName),
		errors.Is(err, domain.ErrEmptyProductSlug),
		errors.Is(err, domain.ErrEmptyCategoryID),
		errors.Is(err, domain.ErrCategoryNotFound), // category_id noto'g'ri/mavjud emas - so'rov xatosi
		errors.Is(err, domain.ErrTooManyProductImages),
		errors.Is(err, domain.ErrInvalidRating),
		errors.Is(err, domain.ErrNegativeStock),
		errors.Is(err, domain.ErrNegativeSoldCount),
		errors.Is(err, domain.ErrDiscountTooHigh),
		errors.Is(err, domain.ErrNegativeAmount),
		errors.Is(err, domain.ErrInvalidCurrency),
		errors.Is(err, domain.ErrCurrencyMismatch),
		errors.Is(err, media.ErrEmptyFile),
		errors.Is(err, media.ErrFileTooLarge),
		errors.Is(err, media.ErrUnsupportedType):
		response.Error(w, http.StatusBadRequest, err.Error()) // 400

	case errors.Is(err, domain.ErrProductSlugTaken):
		response.Error(w, http.StatusConflict, err.Error()) // 409

	case errors.Is(err, domain.ErrProductNotFound):
		response.Error(w, http.StatusNotFound, err.Error()) // 404

	default:
		log.Println(err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
