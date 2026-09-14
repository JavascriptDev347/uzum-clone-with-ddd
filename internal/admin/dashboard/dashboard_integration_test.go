package dashboard

// Bu paket ataylab bounded context emas (domain/application/infrastructure ajratilmagan),
// shuning uchun boshqa hech qanday joyda ishlatilgan soxta repository/uploader naqshlari bu
// yerga to'g'ri kelmaydi - agregatsiyalar xom SQL bo'lgani uchun yagona ma'noli test haqiqiy
// Postgres'ga qarshi ishlaydigan integratsion testdir. Shu sabab bu fayl white-box
// (`package dashboard`, _test emas) - getSummary/getRevenueHistory/getLowStock kabi
// eksport qilinmagan so'rov funksiyalarini to'g'ridan-to'g'ri chaqiradi.
//
// Test haqiqiy Postgres'ga ulana olmasa (`make up && make migrate-up` oldindan ishga
// tushirilmagan bo'lsa), avtomatik ravishda skip qilinadi - `go test ./...` boshqa hech bir
// joyda DB talab qilmaydi, shuning uchun bu paket ham hech kimning muhitini buzmaydi.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	catalogdomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	catalogpg "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/infrastructure/postgres"
	orderingdomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	orderingpg "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/infrastructure/postgres"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/money"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// testDB - haqiqiy Postgres'ga ulanadi. Ulanib bo'lmasa testni skip qiladi.
func testDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := database.NewPostgresDB(testDSN(t))
	if err != nil {
		t.Skipf("postgres bilan bog'lanib bo'lmadi (%v) - avval `make up && make migrate-up` ishga tushiring", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// testDSN - DSN'ni OS environment o'zgaruvchilaridan (CI'da bo'lgani kabi) yoki bo'lmasa
// repo ildizidagi ".env" faylidan yig'adi. pkg/config'dagi Load() bu yerda ishlatilmaydi -
// u ".env"ni joriy ish papkasiga nisbatan qidiradi, `go test` esa paket papkasidan ishga
// tushadi (repo ildizidan emas), shuning uchun bu yordamchi runtime.Caller orqali repo
// ildiziga nisbatan aniq yo'l bilan o'qiydi.
func testDSN(t *testing.T) string {
	t.Helper()

	fileVals := loadRepoRootEnv(t)
	get := func(key, def string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		if v, ok := fileVals[key]; ok && v != "" {
			return v
		}
		return def
	}

	host := get("DB_HOST", "localhost")
	port := get("DB_PORT", "5432")
	user := get("DB_USER", "postgres")
	password := get("DB_PASSWORD", "")
	name := get("DB_NAME", "uzum_clone")
	sslmode := get("DB_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
}

func loadRepoRootEnv(t *testing.T) map[string]string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}
	// internal/admin/dashboard/ -> repo ildizi uchta katalog yuqorida.
	envPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", ".env")

	data, err := os.ReadFile(envPath)
	if err != nil {
		return nil
	}

	values := make(map[string]string)
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return values
}

// seedCategory - test uchun bitta kategoriya yaratadi va uni test oxirida o'chiradi.
func seedCategory(t *testing.T, db *sqlx.DB) string {
	t.Helper()

	repo := catalogpg.NewPostgresCategoryRepository(db)
	id := uuid.New().String()
	category, err := catalogdomain.NewCategory(id, "Test toifa "+id, "Test category "+id, "Тест категория "+id, "", "")
	if err != nil {
		t.Fatalf("unexpected error building category: %v", err)
	}
	if err := repo.Save(context.Background(), category); err != nil {
		t.Fatalf("unexpected error saving category: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	})
	return id
}

// seedProduct - test uchun bitta mahsulot yaratadi (berilgan narx va zaxira bilan) va uni
// test oxirida o'chiradi.
func seedProduct(t *testing.T, db *sqlx.DB, categoryID string, priceAmount float64, stock int) *catalogdomain.Product {
	t.Helper()

	repo := catalogpg.NewPostgresProductRepository(db)
	id := uuid.New().String()
	price, err := money.NewMoney(priceAmount, "UZS")
	if err != nil {
		t.Fatalf("unexpected error building money: %v", err)
	}
	product, err := catalogdomain.NewProduct(catalogdomain.NewProductParams{
		ID:         id,
		NameUz:     "Test mahsulot " + id,
		NameEng:    "Test product " + id,
		NameRu:     "Тест продукт " + id,
		CategoryID: categoryID,
		Price:      price,
		Slug:       "test-product-" + id,
		Stock:      stock,
	})
	if err != nil {
		t.Fatalf("unexpected error building product: %v", err)
	}
	if err := repo.Save(context.Background(), product); err != nil {
		t.Fatalf("unexpected error saving product: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM products WHERE id = $1`, id)
	})
	return product
}

// seedOrder - berilgan itemlar, to'lov holati va yetkazib berish holati bilan bitta
// buyurtma yaratadi va uni test oxirida o'chiradi (order_items ON DELETE CASCADE orqali
// birga o'chadi).
func seedOrder(t *testing.T, db *sqlx.DB, items []orderingdomain.OrderItem, paid bool, deliveryStatus orderingdomain.DeliveryStatus) *orderingdomain.Order {
	t.Helper()

	repo := orderingpg.NewPostgresOrderRepository(db)

	info, err := orderingdomain.NewCustomerInfo("Tashkent, Test 1", "+998901234567", "")
	if err != nil {
		t.Fatalf("unexpected error building customer info: %v", err)
	}

	order, err := orderingdomain.NewOrder(nil, info, items)
	if err != nil {
		t.Fatalf("unexpected error building order: %v", err)
	}
	if paid {
		order.MarkPaid()
	}
	switch deliveryStatus {
	case orderingdomain.DeliveryStatusPreparing:
		// standart holat - hech narsa qilish shart emas.
	case orderingdomain.DeliveryStatusCancelled:
		if err := order.Cancel(); err != nil {
			t.Fatalf("unexpected error cancelling order: %v", err)
		}
	default:
		if err := order.AdvanceDelivery(deliveryStatus); err != nil {
			t.Fatalf("unexpected error advancing delivery status: %v", err)
		}
	}

	if err := repo.Save(context.Background(), order); err != nil {
		t.Fatalf("unexpected error saving order: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.Exec(`DELETE FROM orders WHERE id = $1`, order.ID())
	})
	return order
}

func mustOrderItem(t *testing.T, product *catalogdomain.Product, quantity int) orderingdomain.OrderItem {
	t.Helper()
	item, err := orderingdomain.NewOrderItem(product.ID(), product.NameUz(), product.Price(), quantity)
	if err != nil {
		t.Fatalf("unexpected error building order item: %v", err)
	}
	return item
}

// TestGetSummary_Integration - eng muhim biznes qoidani sinaydi: to'lanmagan buyurtmalar va
// to'langan-lekin-bekor-qilingan buyurtmalar daromadga (yoki sotilgan birliklarga) qo'shilmasligi
// kerak. Mavjud (test bilan bog'liq bo'lmagan) ma'lumotlarga chidamli bo'lish uchun mutlaq
// qiymatlar emas, DELTA'lar tekshiriladi.
func TestGetSummary_Integration(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	before, err := getSummary(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	categoryID := seedCategory(t, db)
	productA := seedProduct(t, db, categoryID, 1000, 10) // 1000 so'm
	productB := seedProduct(t, db, categoryID, 500, 10)  // 500 so'm

	// Hisobga olinishi kerak: to'langan + Delivered, 2*1000 + 1*500 = 2500, 3 birlik.
	seedOrder(t, db,
		[]orderingdomain.OrderItem{mustOrderItem(t, productA, 2), mustOrderItem(t, productB, 1)},
		true, orderingdomain.DeliveryStatusDelivered,
	)
	// Hisobga olinmasligi kerak: to'lanmagan.
	seedOrder(t, db, []orderingdomain.OrderItem{mustOrderItem(t, productA, 5)}, false, orderingdomain.DeliveryStatusPreparing)
	// Hisobga olinmasligi kerak: to'langan, lekin bekor qilingan - bu aynan tekshirilishi
	// so'ralgan qoida.
	seedOrder(t, db, []orderingdomain.OrderItem{mustOrderItem(t, productA, 7)}, true, orderingdomain.DeliveryStatusCancelled)

	after, err := getSummary(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := after.TotalProducts - before.TotalProducts; got != 2 {
		t.Fatalf("expected total_products delta 2, got %d", got)
	}
	if got := after.TotalUnitsSold - before.TotalUnitsSold; got != 3 {
		t.Fatalf("expected total_units_sold delta 3 (unpaid and cancelled-but-paid orders excluded), got %d", got)
	}
	if got := after.TotalRevenue - before.TotalRevenue; got != 2500 {
		t.Fatalf("expected total_revenue delta 2500 (unpaid and cancelled-but-paid orders excluded), got %v", got)
	}
}

// TestGetRevenueHistory_Integration - bugungi kun uchun daromad "day" davri bo'yicha to'g'ri
// guruhlanishini tekshiradi (delta orqali, mavjud ma'lumotlarga chidamli).
func TestGetRevenueHistory_Integration(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	today := time.Now().Truncate(24 * time.Hour)

	revenueForToday := func(points []RevenuePoint) float64 {
		for _, p := range points {
			if p.Period.Truncate(24 * time.Hour).Equal(today) {
				return p.Revenue
			}
		}
		return 0
	}

	before, err := getRevenueHistory(ctx, db, "day")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	beforeToday := revenueForToday(before)

	categoryID := seedCategory(t, db)
	product := seedProduct(t, db, categoryID, 2000, 10)
	seedOrder(t, db, []orderingdomain.OrderItem{mustOrderItem(t, product, 3)}, true, orderingdomain.DeliveryStatusHandedToCourier)

	after, err := getRevenueHistory(ctx, db, "day")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	afterToday := revenueForToday(after)

	if got := afterToday - beforeToday; got != 6000 {
		t.Fatalf("expected today's revenue bucket to grow by 6000 (2000*3), got %v", got)
	}
}

func TestGetRevenueHistory_Integration_RejectsInvalidPeriodAtHandlerLevel(t *testing.T) {
	if !revenueHistoryPeriods["day"] || !revenueHistoryPeriods["month"] {
		t.Fatal("expected 'day' and 'month' to be the supported periods")
	}
	if revenueHistoryPeriods["year"] {
		t.Fatal("expected 'year' not to be a supported period")
	}
}

// TestGetLowStock_Integration - eng past zaxirali mahsulot ro'yxatda ekanligini va
// o'chirilgan (soft-deleted) mahsulotlar butunlay chiqarib tashlanishini tekshiradi.
// Mavjud (test bilan bog'liq bo'lmagan) ma'lumotlar bilan tortishmaslik uchun zaxira
// qasddan haqiqatda bo'lishi mumkin bo'lgan eng past qiymatdan past qilib majburlanadi
// (domen konstruktori manfiy zaxirani rad etadi, lekin DB darajasida bunday cheklov yo'q -
// bu faqat test uchun xom SQL orqali qilinadi, productning o'zi domen orqali to'g'ri
// yaratilgandan keyin).
func TestGetLowStock_Integration(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	categoryID := seedCategory(t, db)

	lowest := seedProduct(t, db, categoryID, 1000, 0)
	if _, err := db.Exec(`UPDATE products SET stock = $1 WHERE id = $2`, -1_000_000, lowest.ID()); err != nil {
		t.Fatalf("unexpected error forcing stock to a sentinel minimum: %v", err)
	}

	deletedButLower := seedProduct(t, db, categoryID, 1000, 0)
	if _, err := db.Exec(`UPDATE products SET stock = $1 WHERE id = $2`, -2_000_000, deletedButLower.ID()); err != nil {
		t.Fatalf("unexpected error forcing stock to a sentinel minimum: %v", err)
	}
	repo := catalogpg.NewPostgresProductRepository(db)
	if err := repo.SoftDelete(ctx, deletedButLower.ID()); err != nil {
		t.Fatalf("unexpected error soft-deleting product: %v", err)
	}

	products, err := getLowStock(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) > lowStockLimit {
		t.Fatalf("expected at most %d products, got %d", lowStockLimit, len(products))
	}

	var found bool
	for i, p := range products {
		if p.ID == deletedButLower.ID() {
			t.Fatal("expected soft-deleted product to be excluded from low-stock results")
		}
		if p.ID == lowest.ID() {
			found = true
			if i != 0 {
				t.Fatalf("expected the sentinel-lowest product to rank first, got index %d", i)
			}
		}
		if i > 0 && p.Stock < products[i-1].Stock {
			t.Fatalf("expected results sorted ascending by stock, got %d after %d", p.Stock, products[i-1].Stock)
		}
	}
	if !found {
		t.Fatal("expected the sentinel-lowest product to appear in the low-stock results")
	}
}
