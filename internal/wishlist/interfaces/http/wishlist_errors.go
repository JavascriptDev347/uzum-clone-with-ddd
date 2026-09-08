package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

var ErrUnauthorized = errors.New("wishlist: unauthorized")

func writeWishlistError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case errors.Is(err, domain.ErrEmptyProductID),
		errors.Is(err, domain.ErrEmptyUserID),
		errors.Is(err, domain.ErrEmptyWishlistID):
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.Envelope{Error: err.Error()})

	case errors.Is(err, domain.ErrProductAlreadyInWishlist):
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response.Envelope{Error: err.Error()})

	case errors.Is(err, domain.ErrProductNotInWishlist),
		errors.Is(err, domain.ErrWishlistNotFound):
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.Envelope{Error: err.Error()})

	case errors.Is(err, ErrUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.Envelope{Error: "unauthorized"})

	default:
		log.Printf("wishlist: kutilmagan xato: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.Envelope{Error: "internal server error"})
	}
}
