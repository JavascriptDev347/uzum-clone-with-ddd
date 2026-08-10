package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

type ProductHandler struct {
	createUc CreateUseCase
}

func NewProductHandler(createUc CreateUseCase) *ProductHandler {
	return &ProductHandler{
		createUc: createUc,
	}
}

// CreateProduct godoc
//
//	@Summary		Yangi mahsulot yaratish
//	@Description	Nomi, narxi, valyutasi va kategoriyasi orqali yangi mahsulot yaratadi
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateProductRequest	true	"Mahsulot ma'lumotlari"
//	@Success		201		{object}	response.Envelope{data=CreateProductResponse}	"Mahsulot yaratildi"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov tanasi yoki validatsiya xatosi"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.createUc.Execute(r.Context(), application.CreateProductInput{
		Name:       req.Name,
		Amount:     req.Amount,
		Currency:   req.Currency,
		CategoryID: req.CategoryID,
	})
	if err != nil {
		writeCatalogError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, CreateProductResponse{
		ID:         output.ID,
		Name:       output.Name,
		CategoryID: output.CategoryID,
		Amount:     output.PriceAmount,   // <-- Output field nomi boshqacha
		Currency:   output.PriceCurrency, // <-- shu yerga e'tibor
		CreatedAt:  output.CreatedAt.Format(time.RFC3339),
	})
}
