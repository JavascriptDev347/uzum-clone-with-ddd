package dashboard

import (
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/jmoiron/sqlx"
)

// SummaryHandler godoc
//
//	@Summary		Umumiy statistika
//	@Description	Jami mahsulotlar soni (o'chirilganlardan tashqari), jami sotilgan birliklar va jami daromad (faqat to'langan va bekor qilinmagan buyurtmalar bo'yicha). Faqat admin uchun.
//	@Tags			admin-dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Envelope{data=SummaryResponse}	"Umumiy statistika"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403	{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/dashboard/summary [get]
func SummaryHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summary, err := getSummary(r.Context(), db)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		response.Success(w, http.StatusOK, summary)
	}
}

// RevenueHistoryHandler godoc
//
//	@Summary		Daromad tarixi
//	@Description	Daromadni davr (kun yoki oy) bo'yicha guruhlab qaytaradi, faqat to'langan va bekor qilinmagan buyurtmalar bo'yicha - frontend'dagi oddiy grafik uchun. Faqat admin uchun.
//	@Tags			admin-dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Param			period	query		string	true	"Guruhlash davri: day yoki month"
//	@Success		200		{object}	response.Envelope{data=[]RevenuePoint}	"Daromad tarixi"
//	@Failure		400		{object}	response.Envelope	"period 'day' yoki 'month' bo'lishi kerak"
//	@Failure		401		{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403		{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/dashboard/revenue-history [get]
func RevenueHistoryHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if !revenueHistoryPeriods[period] {
			response.Error(w, http.StatusBadRequest, "period 'day' yoki 'month' bo'lishi kerak")
			return
		}

		points, err := getRevenueHistory(r.Context(), db, period)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		response.Success(w, http.StatusOK, points)
	}
}

// LowStockHandler godoc
//
//	@Summary		Kam zaxirali mahsulotlar
//	@Description	Eng kam zaxiraga ega (o'chirilmagan) top-5 mahsulot. Faqat admin uchun.
//	@Tags			admin-dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Envelope{data=[]LowStockProduct}	"Kam zaxirali mahsulotlar"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403	{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/dashboard/low-stock [get]
func LowStockHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := getLowStock(r.Context(), db)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}
		response.Success(w, http.StatusOK, products)
	}
}
