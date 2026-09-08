package http

import (
	"encoding/json"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/wishlist/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
)

type WishlistHandler struct {
	addUseCase    *application.AddToWishlistUseCase
	removeUseCase *application.RemoveFromWishlistUseCase
	getUseCase    *application.GetWishlistUseCase
}

func NewWishlistHandler(
	addUC *application.AddToWishlistUseCase,
	removeUC *application.RemoveFromWishlistUseCase,
	getUC *application.GetWishlistUseCase,
) *WishlistHandler {
	return &WishlistHandler{
		addUseCase:    addUC,
		removeUseCase: removeUC,
		getUseCase:    getUC,
	}
}

// @Summary Add product to wishlist
// @Tags wishlist
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Success 204
// @Failure 400 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Router /wishlist/items/{product_id} [post]
func (h *WishlistHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeWishlistError(w, ErrUnauthorized)
		return
	}

	productID := chi.URLParam(r, "product_id")

	err := h.addUseCase.Execute(r.Context(), application.AddToWishlistInput{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Remove product from wishlist
// @Tags wishlist
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Success 204
// @Failure 404 {object} response.Envelope
// @Router /wishlist/items/{product_id} [delete]
func (h *WishlistHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeWishlistError(w, ErrUnauthorized)
		return
	}

	productID := chi.URLParam(r, "product_id")

	err := h.removeUseCase.Execute(r.Context(), application.RemoveFromWishlistInput{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get current user's wishlist
// @Tags wishlist
// @Security BearerAuth
// @Success 200 {object} response.Envelope{data=application.GetWishlistOutput}
// @Router /wishlist [get]
func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeWishlistError(w, ErrUnauthorized)
		return
	}

	output, err := h.getUseCase.Execute(r.Context(), userID)
	if err != nil {
		writeWishlistError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.Envelope{Data: output})
}
