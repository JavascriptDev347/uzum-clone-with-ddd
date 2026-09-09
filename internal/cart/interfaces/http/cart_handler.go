package http

import (
	"encoding/json"
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/cart/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CartHandler struct {
	addItemUseCase            *application.AddItemUseCase
	removeItemUseCase         *application.RemoveItemUseCase
	updateItemQuantityUseCase *application.UpdateItemQuantityUseCase
	getCartUseCase            *application.GetCartUseCase
}

func NewCartHandler(
	addItemUC *application.AddItemUseCase,
	removeItemUC *application.RemoveItemUseCase,
	updateItemQuantityUC *application.UpdateItemQuantityUseCase,
	getCartUC *application.GetCartUseCase,
) *CartHandler {
	return &CartHandler{
		addItemUseCase:            addItemUC,
		removeItemUseCase:         removeItemUC,
		updateItemQuantityUseCase: updateItemQuantityUC,
		getCartUseCase:            getCartUC,
	}
}

type addItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type updateItemQuantityRequest struct {
	Quantity int `json:"quantity"`
}

// @Summary Get current user's cart
// @Tags cart
// @Security BearerAuth
// @Success 200 {object} response.Envelope{data=application.CartView}
// @Failure 401 {object} response.Envelope
// @Router /cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeCartError(w, ErrUnauthorized)
		return
	}

	view, err := h.getCartUseCase.Execute(r.Context(), userID)
	if err != nil {
		writeCartError(w, err)
		return
	}

	response.Success(w, http.StatusOK, view)
}

// @Summary Add product to cart
// @Tags cart
// @Security BearerAuth
// @Param input body addItemRequest true "Product ID and quantity"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Router /cart/items [post]
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeCartError(w, ErrUnauthorized)
		return
	}

	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.addItemUseCase.Execute(r.Context(), application.AddItemInput{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		writeCartError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Mahsulot savatga qo'shildi")
}

// @Summary Update cart item quantity
// @Tags cart
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param input body updateItemQuantityRequest true "New quantity"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Router /cart/items/{product_id} [put]
func (h *CartHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeCartError(w, ErrUnauthorized)
		return
	}

	productID := chi.URLParam(r, "product_id")

	var req updateItemQuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.updateItemQuantityUseCase.Execute(r.Context(), application.UpdateItemQuantityInput{
		UserID:    userID,
		ProductID: productID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		writeCartError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Savatdagi mahsulot miqdori yangilandi")
}

// @Summary Remove product from cart
// @Tags cart
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Router /cart/items/{product_id} [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeCartError(w, ErrUnauthorized)
		return
	}

	productID := chi.URLParam(r, "product_id")

	if err := h.removeItemUseCase.Execute(r.Context(), userID, productID); err != nil {
		writeCartError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Mahsulot savatdan o'chirildi")
}
