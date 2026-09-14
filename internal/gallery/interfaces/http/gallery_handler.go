package http

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type GalleryHandler struct {
	createUseCase *application.CreateGalleryPostUseCase
	listUseCase   *application.ListGalleryPostsUseCase
	deleteUseCase *application.DeleteGalleryPostUseCase
}

func NewGalleryHandler(
	createUC *application.CreateGalleryPostUseCase,
	listUC *application.ListGalleryPostsUseCase,
	deleteUC *application.DeleteGalleryPostUseCase,
) *GalleryHandler {
	return &GalleryHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
		deleteUseCase: deleteUC,
	}
}

// parsePagination - "page" va "page_size" query parametrlarini o'qiydi, bo'sh/noto'g'ri
// bo'lsa 0 qaytaradi (standart qiymatlar use case darajasida qo'llanadi).
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

// CreateGalleryPost godoc
//
//	@Summary		Yangi galereya posti yaratish
//	@Description	Rasm(lar) (eng ko'pi bilan 3 ta) va ixtiyoriy tavsif bilan yangi galereya posti yaratadi. Faqat admin uchun.
//	@Tags			gallery
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			images		formData	file	false	"Rasmlar (eng ko'pi bilan 3 ta)"
//	@Param			description	formData	string	false	"Tavsif"
//	@Success		201			{object}	response.Envelope{data=GalleryPostResponse}	"Post yaratildi"
//	@Failure		400			{object}	response.Envelope	"Noto'g'ri so'rov, 3 tadan ortiq rasm yoki yaroqsiz fayl"
//	@Failure		401			{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403			{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/gallery [post]
func (h *GalleryHandler) CreateGalleryPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "fayl hajmi juda katta yoki noto'g'ri format")
		return
	}

	imageFiles := r.MultipartForm.File["images"]
	if len(imageFiles) > domain.MaxGalleryImages {
		response.Error(w, http.StatusBadRequest, "eng ko'pi bilan 3 ta rasm yuklash mumkin")
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

		// Xavfsizlik: foydalanuvchi yuklagan asl fayl nomi (header.Filename) hech qachon
		// to'g'ridan-to'g'ri public_id sifatida ishlatilmaydi - faqat kengaytmasi olinadi,
		// nomning o'zi uuid bilan almashtiriladi (Category'da avval tuzatilgan, Product'da
		// hali yo'q muammo - bu yerda boshidanoq to'g'ri qilinadi).
		ext := filepath.Ext(fh.Filename)
		images = append(images, media.UploadInput{
			FileName:    "gallery-images/" + uuid.New().String() + ext,
			ContentType: fh.Header.Get("Content-Type"),
			Data:        data,
		})
	}

	post, err := h.createUseCase.Execute(r.Context(), application.CreateGalleryPostInput{
		Images:      images,
		Description: r.FormValue("description"),
	})
	if err != nil {
		writeGalleryError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, ToGalleryPostResponse(post))
}

// ListGalleryPosts godoc
//
//	@Summary		Galereya postlari
//	@Description	Barcha galereya postlarini sahifalab qaytaradi. Ochiq - autentifikatsiya shart emas.
//	@Tags			gallery
//	@Produce		json
//	@Param			page		query		int	false	"Sahifa raqami (default 1)"
//	@Param			page_size	query		int	false	"Sahifadagi elementlar soni (default 20, max 100)"
//	@Success		200			{object}	response.Envelope{data=response.PaginatedResult}	"Postlar ro'yxati"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/gallery [get]
func (h *GalleryHandler) ListGalleryPosts(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)

	posts, total, err := h.listUseCase.Execute(r.Context(), page, pageSize)
	if err != nil {
		writeGalleryError(w, err)
		return
	}

	page, pageSize = application.NormalizeGalleryPagination(page, pageSize)
	response.Success(w, http.StatusOK, response.NewPaginatedResult(ToGalleryPostResponses(posts), total, page, pageSize))
}

// DeleteGalleryPost godoc
//
//	@Summary		Galereya postini o'chirish
//	@Description	Postni bazadan o'chiradi va uning rasmlarini S3'dan ham tozalaydi. Faqat admin uchun.
//	@Tags			gallery
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	response.Envelope	"Post o'chirildi"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403	{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		404	{object}	response.Envelope	"Post topilmadi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/gallery/{id} [delete]
func (h *GalleryHandler) DeleteGalleryPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.deleteUseCase.Execute(r.Context(), id); err != nil {
		writeGalleryError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "Post o'chirildi")
}
