package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

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

// parseTagList - vergul bilan ajratilgan form qiymatini ("romantic, birthday") tag ro'yxatiga aylantiradi.
func parseTagList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// CreateProduct godoc
//
//		@Summary		Yangi mahsulot yaratish
//		@Description	Gullar do'koni mahsulotini barcha xususiyatlari (narx, chegirma, rasmlar, qadoqlash va h.k.) bilan yaratadi. Faqat admin uchun.
//		@Tags			products
//		@Accept			multipart/form-data
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			name					formData	string	true	"Mahsulot nomi"
//		@Param			description				formData	string	false	"Tavsif"
//		@Param			category_id				formData	string	true	"Kategoriya ID"
//		@Param			amount					formData	int		true	"Narx (tiyin/kopeykada)"
//		@Param			currency				formData	string	true	"Valyuta (masalan: UZS)"
//		@Param			discount_amount			formData	int		false	"Chegirma narxi (ixtiyoriy)"
//		@Param			slug					formData	string	false	"Slug (bo'sh bo'lsa nomidan avtomatik yasaladi)"
//		@Param			video_url_youtube		formData	string	false	"YouTube video URL"
//		@Param			video_url_instagram		formData	string	false	"Instagram video URL"
//		@Param			is_available			formData	boolean	false	"Sotuvda bor yoki yo'qligi (default true)"
//		@Param			rating					formData	number	false	"Reyting 1-5 (default 1)"
//		@Param			stock					formData	int		false	"Ombordagi soni"
//		@Param			flower_types			formData	string	false	"Gullar turi, vergul bilan"
//		@Param			color					formData	string	false	"Rangi"
//		@Param			stem_count				formData	int		false	"Buketdagi gullar soni"
//		@Param			packaging_type			formData	string	true	"Qadoqlash turi: bucket, box yoki vase"
//		@Param			freshness_lifespan		formData	int		false	"Saqlanish muddati (kunlarda, 1-7, default 1)"
//		@Param			care_instructions		formData	string	false	"Parvarish bo'yicha ko'rsatma"
//		@Param			occasions				formData	string	false	"Bayramlar/holatlar (tag), vergul bilan"
//		@Param			compatible_addons		formData	string	false	"Mos qo'shimchalar, vergul bilan"
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

	var stemCount int
	if raw := r.FormValue("stem_count"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "noto'g'ri stem_count qiymati")
			return
		}
		stemCount = v
	}

	var freshnessLifespan int
	if raw := r.FormValue("freshness_lifespan"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "noto'g'ri freshness_lifespan qiymati")
			return
		}
		freshnessLifespan = v
	}

	var careInstructions *string
	if raw := r.FormValue("care_instructions"); raw != "" {
		careInstructions = &raw
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
		Name:              r.FormValue("name"),
		Description:       r.FormValue("description"),
		Images:            images,
		VideoURLYoutube:   r.FormValue("video_url_youtube"),
		VideoURLInstagram: r.FormValue("video_url_instagram"),
		CategoryID:        r.FormValue("category_id"),
		Amount:            amount,
		Currency:          r.FormValue("currency"),
		DiscountAmount:    discountAmount,
		Slug:              r.FormValue("slug"),
		IsAvailable:       isAvailable,
		Rating:            rating,
		Stock:             stock,
		FlowerTypes:       parseTagList(r.FormValue("flower_types")),
		Color:             r.FormValue("color"),
		StemCount:         stemCount,
		PackagingType:     r.FormValue("packaging_type"),
		FreshnessLifespan: freshnessLifespan,
		CareInstructions:  careInstructions,
		Occasions:         parseTagList(r.FormValue("occasions")),
		CompatibleAddons:  parseTagList(r.FormValue("compatible_addons")),
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
//	@Description	Faol (o'chirilmagan) mahsulotlar ro'yxati, nomi va/yoki category_id bo'yicha filtrlash mumkin
//	@Tags			products
//	@Produce		json
//	@Param			search		query		string	false	"Mahsulot nomi bo'yicha qidirish"
//	@Param			category_id	query		string	false	"Kategoriya ID bo'yicha filtrlash"
//	@Success		200			{object}	response.Envelope{data=[]application.ProductOutput}	"Mahsulotlar"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category_id")
	products, err := h.getUc.Execute(r.Context(), search, categoryID)
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, products)
}

// GetProductsByCategory godoc
//
//	@Summary		Kategoriya bo'yicha mahsulotlarni olish
//	@Description	Berilgan category_id ga tegishli faol (o'chirilmagan) mahsulotlar ro'yxati
//	@Tags			products
//	@Produce		json
//	@Param			id		path		string	true	"Kategoriya ID"
//	@Param			search	query		string	false	"Mahsulot nomi bo'yicha qidirish"
//	@Success		200		{object}	response.Envelope{data=[]application.ProductOutput}	"Mahsulotlar"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/categories/{id}/products [get]
func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	search := r.URL.Query().Get("search")
	products, err := h.getUc.Execute(r.Context(), search, categoryID)
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, products)
}

// GetProduct godoc
//
//	@Summary		Mahsulotni olish
//	@Description	Mahsulotni ID bo'yicha olish
//	@Tags			products
//	@Produce		json
//	@Param			id	path		string	true	"Mahsulot ID"
//	@Success		200	{object}	response.Envelope{data=application.ProductOutput}	"Mahsulot"
//	@Failure		404	{object}	response.Envelope	"Mahsulot topilmadi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.getByIdUc.Execute(r.Context(), id)
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
//	@Success		200		{object}	response.Envelope{data=application.ProductOutput}	"Mahsulot"
//	@Failure		404		{object}	response.Envelope	"Mahsulot topilmadi"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	product, err := h.getBySlug.Execute(r.Context(), slug)
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
//		@Description	Barcha mahsulotlarni (o'chirilganlarni ham) qaytaradi
//		@Tags			products
//	 @Security		BearerAuth
//		@Produce		json
//		@Param			search		query		string	false	"Mahsulot nomi bo'yicha qidirish"
//		@Param			category_id	query		string	false	"Kategoriya ID bo'yicha filtrlash"
//		@Success		200			{object}	response.Envelope{data=[]application.ProductOutputForAdmin}	"Mahsulotlar"
//		@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//		@Router			/products/admin [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	categoryID := r.URL.Query().Get("category_id")
	products, err := h.getAllUc.Execute(r.Context(), search, categoryID)
	if err != nil {
		writeProductError(w, err)
		return
	}

	response.Success(w, http.StatusOK, products)
}
