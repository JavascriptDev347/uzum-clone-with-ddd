package http

import (
	"errors"
	"net/http"

	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

var ErrUnauthorized = errors.New("ordering: unauthorized")

func writeOrderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyAddress):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidPhone):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmptyOrderItems):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmptyOrderItemProductID):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrEmptyOrderItemName):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidQuantity):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInsufficientStock):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrInvalidPaymentStatus):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidDeliveryStatus):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrInvalidDeliveryTransition):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrOrderAlreadyShipped):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrOrderNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrOrderAccessDenied):
		response.Error(w, http.StatusForbidden, err.Error())
	case errors.Is(err, catalog.ErrProductNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrUnauthorized):
		response.Error(w, http.StatusUnauthorized, "unauthorized")
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
