package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
)

type CategoryHandler struct {
	createUc  CreateCategoryUseCase
	getUc     GetCategoriesUseCase
	getByIdUc GetCategoryUseCase
	updateUc  UpdateCategoryUseCase
}

func NewCategoryHandler(createUc CreateCategoryUseCase, getUc GetCategoriesUseCase, getByIdUc GetCategoryUseCase, updateUc UpdateCategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		createUc:  createUc,
		getUc:     getUc,
		getByIdUc: getByIdUc,
		updateUc:  updateUc,
	}
}

// CreateCategory godoc
//
//		@Summary		Yangi kategoriya yaratish
//		@Description	Nomi va boshqa ma'lumotlar orqali yangi kategoriya yaratadi
//		@Tags			categories
//		@Accept			json
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			request	body		CreateCategoryRequest	true	"Kategoriya ma'lumotlari"
//		@Success		201		{object}	response.Envelope{data=CreateCategoryResponse}	"Kategoriya yaratildi"//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov tanasi yoki validatsiya xatosi"
//		@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.createUc.Execute(r.Context(), application.CreateCategoryInput{
		Name: req.Name,
	})
	if err != nil {
		writeCategoryError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, CreateCategoryResponse{
		ID:        output.ID,
		Name:      output.Name,
		UpdatedAt: output.UpdatedAt.Format(time.RFC3339),
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	})
}

// GetCategories godoc
//
//	@Summary		Kategoriyalarni olish
//	@Description	Kategoriyalarni nomi bo'yicha qidirish
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string	false	"Kategoriya nomi"
//	@Success		200		{object}	response.Envelope{data=[]CreateCategoryResponse}	"Kategoriyalar olish muvaffaqiyatli"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/categories [get]
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	categories, err := h.getUc.Execute(r.Context(), search)
	if err != nil {
		writeCategoryError(w, err)
		return
	}

	response.Success(w, http.StatusOK, categories)
}

// GetCategory godoc
//
//	@Summary		Kategoriyani olish
//	@Description	Kategoriyani ID bo'yicha olish
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Kategoriya ID"
//	@Success		200		{object}	response.Envelope{data=CreateCategoryResponse}	"Kategoriya olish muvaffaqiyatli"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/categories/{id} [get]
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category, err := h.getByIdUc.Execute(r.Context(), id)
	if err != nil {
		writeCategoryError(w, err)
		return
	}

	response.Success(w, http.StatusOK, category)
}

// UpdateCategory godoc
//
//		@Summary		Kategoriyani yangilash
//		@Description	Kategoriyani ID bo'yicha yangilash
//		@Tags			categories
//		@Accept			json
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			id	path		string	true	"Kategoriya ID"
//		@Param			category	body		application.UpdateCategoryInput	true	"Yangilash uchun kategoriya ma'lumotlari"
//		@Success		200		{object}	response.Envelope	"Kategoriya yangilash muvaffaqiyatli"
//		@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input application.UpdateCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeCategoryError(w, err)
		return
	}
	input.ID = id
	if err := h.updateUc.Execute(r.Context(), input); err != nil {
		writeCategoryError(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil)
}
