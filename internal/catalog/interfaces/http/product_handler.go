package http

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/google/uuid"
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
//	@Description	Nomi, narxi, valyutasi, kategoriyasi va rasmi orqali yangi mahsulot yaratadi. Faqat admin huquqiga ega foydalanuvchi uchun.
//	@Tags			products
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		formData	string	true	"Mahsulot nomi"
//	@Param			amount		formData	int		true	"Narx (tiyin/kopeykada)"
//	@Param			currency	formData	string	true	"Valyuta (masalan: UZS)"
//	@Param			categoryId	formData	string	true	"Kategoriya ID"
//	@Param			image		formData	file	false	"Mahsulot rasmi"
//	@Success		201			{object}	response.Envelope{data=CreateProductResponse}	"Mahsulot yaratildi"
//	@Failure		400			{object}	response.Envelope	"Noto'g'ri so'rov tanasi yoki validatsiya xatosi"
//	@Failure		401			{object}	response.Envelope	"Autentifikatsiyadan o'tilmagan"
//	@Failure		403			{object}	response.Envelope	"Huquq yetarli emas"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	var imageData []byte
	var imageContentType, imageFileName string

	file, header, err := r.FormFile("image")
	if err == nil { // rasm ixtiyoriy - yuborilmasa xato emas
		defer file.Close()
		imageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read image", http.StatusBadRequest)
			return
		}
		imageContentType = header.Header.Get("Content-Type")
		imageFileName = "product-images/" + uuid.NewString()
	}

	output, err := h.createUc.Execute(r.Context(), application.CreateProductInput{
		Name:             r.FormValue("name"),
		Amount:           amount,
		Currency:         r.FormValue("currency"),
		CategoryID:       r.FormValue("categoryId"),
		ImageData:        imageData,
		ImageContentType: imageContentType,
		ImageFileName:    imageFileName,
	})
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, CreateProductResponse{
		ID:         output.ID,
		Name:       output.Name,
		CategoryID: output.CategoryID,
		Amount:     output.PriceAmount,
		Currency:   output.PriceCurrency,
		ImageURL:   output.ImageURL,
		CreatedAt:  output.CreatedAt.Format(time.RFC3339),
	})
}
