package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProductHandler struct {
	createUc  CreateProductUseCase
	getUc     GetProductsUseCase
	getByIdUc GetProductUseCase
	getBySlug GetProductBySlugUseCase
	updateUc  UpdateProductUseCase
	deleteUc  DeleteProductUseCase
	getAllUc  GetAllProductsIncludingDeletedUseCase
}

func NewProductHandler(
	createUc CreateProductUseCase,
	getUc GetProductsUseCase,
	getByIdUc GetProductUseCase,
	getBySlug GetProductBySlugUseCase,
	updateUc UpdateProductUseCase,
	deleteUc DeleteProductUseCase,
	getAllUc GetAllProductsIncludingDeletedUseCase,
) *ProductHandler {
	return &ProductHandler{
		createUc:  createUc,
		getUc:     getUc,
		getByIdUc: getByIdUc,
		getBySlug: getBySlug,
		updateUc:  updateUc,
		deleteUc:  deleteUc,
		getAllUc:  getAllUc,
	}
}

// parsePagination - "page" va "page_size" query parametrlarini o'qiydi, bo'sh/noto'g'ri bo'lsa 0 qaytaradi
// (standart qiymatlar use case darajasida qo'llanadi).
func parsePagination(r *http.Request) (page int, pageSize int) {
	if raw := r.URL.Query().Get("page"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			page = v
		}
	}
	if raw := r.URL.Query().Get("page_size"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			pageSize = v
		}
	}
	return page, pageSize
}

// parseLang - "lang" query parametrini o'qiydi (uz/eng/ru), bo'sh yoki noto'g'ri bo'lsa uz qaytadi.
func parseLang(r *http.Request) application.Lang {
	return application.ParseLang(r.URL.Query().Get("lang"))
}

// CreateProduct godoc
//
//		@Summary		Yangi mahsulot yaratish
//		@Description	Yangi mahsulotni barcha xususiyatlari (nomi/tavsifi 3 tilda, narx, chegirma, rasmlar) bilan yaratadi. Faqat admin uchun.
//		@Tags			products
//		@Accept			multipart/form-data
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			name_uz					formData	string	true	"Mahsulot nomi (o'zbekcha)"
//		@Param			name_eng				formData	string	true	"Mahsulot nomi (inglizcha)"
//		@Param			name_ru					formData	string	true	"Mahsulot nomi (ruscha)"
//		@Param			description_uz			formData	string	false	"Tavsif (o'zbekcha)"
//		@Param			description_eng			formData	string	false	"Tavsif (inglizcha)"
//		@Param			description_ru			formData	string	false	"Tavsif (ruscha)"
//		@Param			category_id				formData	string	true	"Kategoriya ID"
//		@Param			amount					formData	int		true	"Narx (so'mda, butun son)"
//		@Param			currency				formData	string	true	"Valyuta (masalan: UZS)"
//		@Param			discount_amount			formData	int		false	"Chegirma narxi, so'mda (ixtiyoriy)"
//		@Param			slug					formData	string	false	"Slug (bo'sh bo'lsa nomidan avtomatik yasaladi)"
//		@Param			is_available			formData	boolean	false	"Sotuvda bor yoki yo'qligi (default true)"
//		@Param			rating					formData	number	false	"Reyting 1-5 (default 1)"
//		@Param			stock					formData	int		false	"Ombordagi soni"
//		@Param			tag_uz					formData	string	false	"Belgi/badge, masalan bestseller (o'zbekcha, ixtiyoriy)"
//		@Param			tag_eng					formData	string	false	"Belgi/badge (inglizcha, ixtiyoriy)"
//		@Param			tag_ru					formData	string	false	"Belgi/badge (ruscha, ixtiyoriy)"
//		@Param			images					formData	file	false	"Mahsulot rasmlari (eng ko'pi bilan 5 ta)"
//		@Success		201			{object}	response.Envelope{data=application.ProductOutput}	"Mahsulot yaratildi"
//		@Failure		400			{object}	response.Envelope	"Noto'g'ri so'rov tanasi yoki validatsiya xatosi"
//		@Failure		401			{object}	response.Envelope	"Autentifikatsiyadan o'tilmagan"
//		@Failure		403			{object}	response.Envelope	"Huquq yetarli emas"
//		@Failure		409			{object}	response.Envelope	"Slug band"
//		@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "fayl hajmi juda katta yoki noto'g'ri format")
		return
	}

	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "noto'g'ri narx (amount)")
		return
	}

	var discountAmount *int64
	if raw := r.FormValue("discount_amount"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "noto'g'ri chegirma narxi (discount_amount)")
			return
		}
		discountAmount = &v
	}

	isAvailable := true
	if raw := r.FormValue("is_available"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "noto'g'ri is_available qiymati")
			return
		}
		isAvailable = v
	}

	var rating float64
	if raw := r.FormValue("rating"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "noto'g'ri rating qiymati")
			return
		}
		rating = v
	}

	var stock int
	if raw := r.FormValue("stock"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "noto'g'ri stock qiymati")
			return
		}
		stock = v
	}

	var tagUz *string
	if raw := r.FormValue("tag_uz"); raw != "" {
		tagUz = &raw
	}
	var tagEng *string
	if raw := r.FormValue("tag_eng"); raw != "" {
		tagEng = &raw
	}
	var tagRu *string
	if raw := r.FormValue("tag_ru"); raw != "" {
		tagRu = &raw
	}

	imageFiles := r.MultipartForm.File["images"]
	if len(imageFiles) > domain.MaxProductImages {
		response.Error(w, http.StatusBadRequest, "eng ko'pi bilan 5 ta rasm yuklash mumkin")
		return
	}

	images := make([]media.UploadInput, 0, len(imageFiles))
	for _, fh := range imageFiles {
		file, err := fh.Open()
		if err != nil {
			response.Error(w, http.StatusBadRequest, "rasmni ochib bo'lmadi")
			return
		}
		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "rasmni o'qishda xatolik")
			return
		}
		images = append(images, media.UploadInput{
			FileName:    "product-images/" + uuid.NewString(),
			ContentType: fh.Header.Get("Content-Type"),
			Data:        data,
		})
	}

	input := application.CreateProductInput{
		NameUz:         r.FormValue("name_uz"),
		NameEng:        r.FormValue("name_eng"),
		NameRu:         r.FormValue("name_ru"),
		DescriptionUz:  r.FormValue("description_uz"),
		DescriptionEng: r.FormValue("description_eng"),
		DescriptionRu:  r.FormValue("description_ru"),
		Images:         images,
		CategoryID:     r.FormValue("category_id"),
		Amount:         amount,
		Currency:       r.FormValue("currency"),
		DiscountAmount: discountAmount,
		Slug:           r.FormValue("slug"),
		IsAvailable:    isAvailable,
		Rating:         rating,
		Stock:          stock,
		TagUz:          tagUz,
		TagEng:         tagEng,
		TagRu:          tagRu,
	}

	output, err := h.createUc.Execute(r.Context(), input)
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, output)
}

// GetProducts godoc
//
//	@Summary		Mahsulotlarni olish
//	@Description	Faol (o'chirilmagan) mahsulotlar ro'yxati, nomi va/yoki category_id bo'yicha filtrlash mumkin. lang bo'yicha localized javob qaytadi.
//	@Tags			products
//	@Produce		json
//	@Param			search		query		string	false	"Mahsulot nomi bo'yicha qidirish"
//	@Param			category_id	query		string	false	"Kategoriya ID bo'yicha filtrlash"
//	@Param			lang		query		string	false	"Til: uz (default), eng yoki ru"
//	@Param			page		query		int		false	"Sahifa raqami (default 1)"
//	@Param			page_size	query		int		false	"Sahifadagi elementlar soni (default 20, max 100)"
//	@Success		200			{object}	response.Envelope{data=response.PaginatedResult}	"Mahsulotlar"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category_id")
	page, pageSize := parsePagination(r)
	lang := parseLang(r)
	products, total, err := h.getUc.Execute(r.Context(), search, categoryID, page, pageSize, lang)
	if err != nil {
		writeProductError(w, err)
		return
	}

	page, pageSize = application.NormalizeProductPagination(page, pageSize)
	response.Success(w, http.StatusOK, response.NewPaginatedResult(products, total, page, pageSize))
}

// GetProductsByCategory godoc
//
//	@Summary		Kategoriya bo'yicha mahsulotlarni olish
//	@Description	Berilgan category_id ga tegishli faol (o'chirilmagan) mahsulotlar ro'yxati
//	@Tags			products
//	@Produce		json
//	@Param			id			path		string	true	"Kategoriya ID"
//	@Param			search		query		string	false	"Mahsulot nomi bo'yicha qidirish"
//	@Param			lang		query		string	false	"Til: uz (default), eng yoki ru"
//	@Param			page		query		int		false	"Sahifa raqami (default 1)"
//	@Param			page_size	query		int		false	"Sahifadagi elementlar soni (default 20, max 100)"
//	@Success		200			{object}	response.Envelope{data=response.PaginatedResult}	"Mahsulotlar"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/categories/{id}/products [get]
func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	search := r.URL.Query().Get("search")
	page, pageSize := parsePagination(r)
	lang := parseLang(r)
	products, total, err := h.getUc.Execute(r.Context(), search, categoryID, page, pageSize, lang)
	if err != nil {
		writeProductError(w, err)
		return
	}

	page, pageSize = application.NormalizeProductPagination(page, pageSize)
	response.Success(w, http.StatusOK, response.NewPaginatedResult(products, total, page, pageSize))
}

// GetProduct godoc
//
//	@Summary		Mahsulotni olish
//	@Description	Mahsulotni ID bo'yicha olish
//	@Tags			products
//	@Produce		json
//	@Param			id		path		string	true	"Mahsulot ID"
//	@Param			lang	query		string	false	"Til: uz (default), eng yoki ru"
//	@Success		200	{object}	response.Envelope{data=application.ProductOutput}	"Mahsulot"
//	@Failure		404	{object}	response.Envelope	"Mahsulot topilmadi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	lang := parseLang(r)
	product, err := h.getByIdUc.Execute(r.Context(), id, lang)
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// GetProductBySlug godoc
//
//	@Summary		Mahsulotni slug bo'yicha olish
//	@Description	Mahsulot sahifasi uchun SEO-friendly slug bo'yicha olish
//	@Tags			products
//	@Produce		json
//	@Param			slug	path		string	true	"Mahsulot slug"
//	@Param			lang	query		string	false	"Til: uz (default), eng yoki ru"
//	@Success		200		{object}	response.Envelope{data=application.ProductOutput}	"Mahsulot"
//	@Failure		404		{object}	response.Envelope	"Mahsulot topilmadi"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	lang := parseLang(r)
	product, err := h.getBySlug.Execute(r.Context(), slug, lang)
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, product)
}

// UpdateProduct godoc
//
//		@Summary		Mahsulotni yangilash
//		@Description	Mahsulotni ID bo'yicha yangilash (JSON body, rasmlar bu yerda o'zgartirilmaydi)
//		@Tags			products
//		@Accept			json
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			id		path		string							true	"Mahsulot ID"
//		@Param			product	body		application.UpdateProductInput	true	"Yangilash uchun mahsulot ma'lumotlari"
//		@Success		200		{object}	response.Envelope	"Mahsulot yangilash muvaffaqiyatli"
//		@Failure		400		{object}	response.Envelope	"Validatsiya xatosi"
//		@Failure		404		{object}	response.Envelope	"Mahsulot topilmadi"
//		@Failure		409		{object}	response.Envelope	"Slug band"
//		@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input application.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.ID = id

	if err := h.updateUc.Execute(r.Context(), input); err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, nil)
}

// DeleteProduct godoc
//
//		@Summary		Mahsulotni o'chirish
//		@Description	Mahsulotni ID bo'yicha o'chirish (soft delete)
//		@Tags			products
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			id	path		string	true	"Mahsulot ID"
//		@Success		200	{object}	response.Envelope	"Mahsulot o'chirish muvaffaqiyatli"
//		@Failure		404	{object}	response.Envelope	"Mahsulot topilmadi"
//		@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.deleteUc.Execute(r.Context(), id); err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Mahsulot o'chirildi")
}

// GetAllProducts godoc
//
//		@Summary		Mahsulotlarni olish - admin (o'chirilganlar bilan birga)
//		@Description	Barcha mahsulotlarni (o'chirilganlarni ham), 3 ta tildagi to'liq ma'lumot bilan qaytaradi
//		@Tags			products
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			search		query		string	false	"Mahsulot nomi bo'yicha qidirish"
//		@Param			category_id	query		string	false	"Kategoriya ID bo'yicha filtrlash"
//		@Param			page		query		int		false	"Sahifa raqami (default 1)"
//		@Param			page_size	query		int		false	"Sahifadagi elementlar soni (default 20, max 100)"
//		@Success		200			{object}	response.Envelope{data=response.PaginatedResult}	"Mahsulotlar"
//		@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/products/admin [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category_id")
	page, pageSize := parsePagination(r)
	products, total, err := h.getAllUc.Execute(r.Context(), search, categoryID, page, pageSize)
	if err != nil {
		writeProductError(w, err)
		return
	}

	page, pageSize = application.NormalizeProductPagination(page, pageSize)
	response.Success(w, http.StatusOK, response.NewPaginatedResult(products, total, page, pageSize))
}
