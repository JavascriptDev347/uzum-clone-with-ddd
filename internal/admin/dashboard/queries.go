// Package dashboard - admin uchun faqat-o'qish uchun agregatsiya qatlami. Bu ATAYLAB
// bounded context EMAS: domain/application/infrastructure ajratilmagan, chunki bu yerda
// himoya qilinadigan hech qanday invariant yo'q - faqat xom SQL agregatsiyalari. Shu sabab
// catalog va ordering'ning jadvallariga to'g'ridan-to'g'ri, repository interfeyssiz murojaat
// qilinadi (boshqa har qanday joyda bunday qilish noto'g'ri bo'lar edi).
package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// SummaryResponse - umumiy statistika: mahsulotlar soni, sotilgan birliklar va daromad.
// Daromad yagona valyutada hisoblanadi deb faraz qilinadi (loyihada amalda har doim shunday) -
// bu qatlam "faqat agregatsiya", valyutalarni solishtirish kabi biznes qoidalarni bilmaydi.
type SummaryResponse struct {
	TotalProducts  int64   `db:"total_products" json:"total_products"`
	TotalUnitsSold int64   `db:"total_units_sold" json:"total_units_sold"`
	TotalRevenue   float64 `db:"total_revenue" json:"total_revenue"`
}

// paidNotCancelledFilter - "to'langan va bekor qilinmagan" buyurtmalar shartini order_items'ni
// orders bilan bog'laydigan har uchala so'rovda takrorlanadi: bekor qilingan (lekin to'langan)
// buyurtmalar daromad sifatida hisoblanmasligi kerak.
const paidNotCancelledFilter = `o.payment_status = 'paid' AND o.delivery_status <> 'cancelled'`

// getSummary - jami mahsulotlar soni (o'chirilganlardan tashqari), jami sotilgan birliklar
// va jami daromad (to'langan va bekor qilinmagan buyurtmalar bo'yicha).
func getSummary(ctx context.Context, db *sqlx.DB) (SummaryResponse, error) {
	var totalProducts int64
	if err := db.GetContext(ctx, &totalProducts, `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`); err != nil {
		return SummaryResponse{}, err
	}

	var totals struct {
		TotalUnitsSold int64   `db:"total_units_sold"`
		TotalRevenue   float64 `db:"total_revenue"`
	}
	query := `
		SELECT
			COALESCE(SUM(oi.quantity), 0) AS total_units_sold,
			COALESCE(SUM(oi.unit_price_amount * oi.quantity), 0) AS total_revenue
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE ` + paidNotCancelledFilter
	if err := db.GetContext(ctx, &totals, query); err != nil {
		return SummaryResponse{}, err
	}

	return SummaryResponse{
		TotalProducts:  totalProducts,
		TotalUnitsSold: totals.TotalUnitsSold,
		TotalRevenue:   totals.TotalRevenue,
	}, nil
}

// RevenuePoint - bitta davr (kun yoki oy) uchun daromad.
type RevenuePoint struct {
	Period  time.Time `db:"period" json:"period"`
	Revenue float64   `db:"revenue" json:"revenue"`
}

// revenueHistoryPeriods - qo'llab-quvvatlanadigan date_trunc birliklari. Handler shu ro'yxat
// bo'yicha tekshiradi - boshqa har qanday qiymat 400 bilan rad etiladi (date_trunc'ning
// birinchi argumenti sifatida ixtiyoriy matn parametrlashtirilishi mumkin bo'lsa-da, faqat
// mantiqan to'g'ri keladigan ikkitasi qabul qilinadi).
var revenueHistoryPeriods = map[string]bool{"day": true, "month": true}

// getRevenueHistory - period bo'yicha guruhlangan daromad tarixi (kunlik yoki oylik), xuddi
// getSummary bilan bir xil "to'langan va bekor qilinmagan" filtri bilan.
func getRevenueHistory(ctx context.Context, db *sqlx.DB, period string) ([]RevenuePoint, error) {
	query := `
		SELECT
			date_trunc(:period, o.created_at) AS period,
			COALESCE(SUM(oi.unit_price_amount * oi.quantity), 0) AS revenue
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE ` + paidNotCancelledFilter + `
		GROUP BY period
		ORDER BY period`

	stmt, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	points := []RevenuePoint{}
	if err := stmt.SelectContext(ctx, &points, map[string]any{"period": period}); err != nil {
		return nil, err
	}
	return points, nil
}

// LowStockProduct - eng kam zaxiraga ega mahsulotlardan biri.
type LowStockProduct struct {
	ID     string `db:"id" json:"id"`
	NameUz string `db:"name_uz" json:"name_uz"`
	Stock  int    `db:"stock" json:"stock"`
}

// lowStockLimit - past zaxirali mahsulotlar ro'yxatida qaytariladigan eng ko'p qatorlar soni.
const lowStockLimit = 5

// getLowStock - eng kam zaxiraga ega (o'chirilmagan) mahsulotlardan top-5.
func getLowStock(ctx context.Context, db *sqlx.DB) ([]LowStockProduct, error) {
	query := fmt.Sprintf(`
		SELECT id, name_uz, stock
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY stock ASC
		LIMIT %d`, lowStockLimit)

	products := []LowStockProduct{}
	if err := db.SelectContext(ctx, &products, query); err != nil {
		return nil, err
	}
	return products, nil
}
