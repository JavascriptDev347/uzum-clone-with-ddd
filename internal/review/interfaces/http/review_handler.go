package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/review/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
)

type ReviewHandler struct {
	submitReviewUseCase      *application.SubmitReviewUseCase
	getProductReviewsUseCase *application.GetProductReviewsUseCase
	getUserReviewsUseCase    *application.GetUserReviewsUseCase
}

func NewReviewHandler(
	submitReviewUC *application.SubmitReviewUseCase,
	getProductReviewsUC *application.GetProductReviewsUseCase,
	getUserReviewsUC *application.GetUserReviewsUseCase,
) *ReviewHandler {
	return &ReviewHandler{
		submitReviewUseCase:      submitReviewUC,
		getProductReviewsUseCase: getProductReviewsUC,
		getUserReviewsUseCase:    getUserReviewsUC,
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

// SubmitReview godoc
//
//	@Summary		Mahsulotga sharh qoldirish
//	@Description	Joriy foydalanuvchi nomidan mahsulotga sharh qoldiradi. Faqat shu mahsulotni sotib olib, yetkazib berilgan (Delivered) buyurtmasi bo'lgan foydalanuvchilar sharh qoldira oladi, va har bir foydalanuvchi bitta mahsulotga faqat bitta marta sharh qoldirishi mumkin.
//	@Tags			reviews
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		SubmitReviewRequest	true	"Sharh ma'lumotlari"
//	@Success		201		{object}	response.Envelope{data=ReviewResponse}	"Sharh yaratildi"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov yoki reyting 1-5 oralig'ida emas"
//	@Failure		401		{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403		{object}	response.Envelope	"Bu mahsulot uchun sharh qoldirish huquqi yo'q"
//	@Failure		409		{object}	response.Envelope	"Bu mahsulot uchun sharh allaqachon qoldirilgan"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/reviews [post]
func (h *ReviewHandler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeReviewError(w, ErrUnauthorized)
		return
	}

	var req SubmitReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	review, err := h.submitReviewUseCase.Execute(r.Context(), application.SubmitReviewInput{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		writeReviewError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, ToReviewResponse(review))
}

// GetProductReviews godoc
//
//	@Summary		Mahsulot sharhlari
//	@Description	Berilgan mahsulot uchun barcha sharhlar ro'yxatini sahifalab qaytaradi. Ochiq - autentifikatsiya shart emas.
//	@Tags			reviews
//	@Produce		json
//	@Param			id			path		string	true	"Mahsulot ID"
//	@Param			page		query		int		false	"Sahifa raqami (default 1)"
//	@Param			page_size	query		int		false	"Sahifadagi elementlar soni (default 20, max 100)"
//	@Success		200			{object}	response.Envelope{data=response.PaginatedResult}	"Sharhlar ro'yxati"
//	@Failure		500			{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/products/{id}/reviews [get]
func (h *ReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	page, pageSize := parsePagination(r)

	reviews, total, err := h.getProductReviewsUseCase.Execute(r.Context(), productID, page, pageSize)
	if err != nil {
		writeReviewError(w, err)
		return
	}

	page, pageSize = application.NormalizeReviewPagination(page, pageSize)
	response.Success(w, http.StatusOK, response.NewPaginatedResult(ToReviewResponses(reviews), total, page, pageSize))
}
