package http

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	createUc  CreateCategoryUseCase
	getUc     GetCategoriesUseCase
	getByIdUc GetCategoryUseCase
	updateUc  UpdateCategoryUseCase
	deleteUc  DeleteCategoryUseCase
}

func NewCategoryHandler(createUc CreateCategoryUseCase, getUc GetCategoriesUseCase, getByIdUc GetCategoryUseCase, updateUc UpdateCategoryUseCase, deleteUc DeleteCategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		createUc:  createUc,
		getUc:     getUc,
		getByIdUc: getByIdUc,
		updateUc:  updateUc,
		deleteUc:  deleteUc,
	}
}

// CreateCategory godoc
//
//		@Summary		Yangi kategoriya yaratish
//		@Description	Nomi va rasm orqali yangi kategoriya yaratadi
//		@Tags			categories
//		@Accept			multipart/form-data
//	    @Security		BearerAuth
//		@Produce		json
//		@Param			name	formData	string	true	"Kategoriya nomi"
//		@Param			image	formData	file	true	"Kategoriya rasmi"
//		@Success		201		{object}	response.Envelope{data=CreateCategoryResponse}	"Kategoriya yaratildi"
//		@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov tanasi yoki validatsiya xatosi"
//		@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "fayl hajmi juda katta yoki noto'g'ri format")
		return
	}

	name := r.FormValue("name")

	file, header, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "rasm yuklanmadi")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "rasmni o'qishda xatolik")
		return
	}

	ext := filepath.Ext(header.Filename) // ".png"
	uniqueFileName := uuid.New().String() + ext
	input := application.CreateCategoryInput{
		Name: name,
		Image: media.UploadInput{
			Data:     data,
			FileName: uniqueFileName,
		},
	}

	output, err := h.createUc.Execute(r.Context(), input)
	if err != nil {
		writeCategoryError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, output)
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

// DeleteCategory godoc
//
//		@Summary		Kategoriyani o'chirish
//		@Description	Kategoriyani ID bo'yicha o'chirish
//		@Tags			categories
//		@Accept			json
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			id	path		string	true	"Kategoriya ID"
//		@Success		200		{object}	response.Envelope	"Kategoriya o'chirish muvaffaqiyatli"
//		@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.deleteUc.Execute(r.Context(), id); err != nil {
		writeCategoryError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Kategoriya o'chirildi")
}
