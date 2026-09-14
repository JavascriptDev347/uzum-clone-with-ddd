package application

import (
	"context"

	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// productQuantity - OrderItem snapshot yaratish uchun kerakli minimal kirish: mahsulot ID va
// miqdor. Checkout (cart'dan) va CreateManualOrder (admin tomonidan qo'lda kiritilgan
// ro'yxatdan) ikkalasi ham shu shaklga tushiriladi, shunda snapshot va zaxira mantig'i
// ikkala use case'da ham takrorlanmaydi.
type productQuantity struct {
	ProductID string
	Quantity  int
}

// buildOrderItemSnapshots - har bir (productID, quantity) jufti uchun Catalog'dan joriy
// mahsulot ma'lumotini o'qiydi va OrderItem'ni quradi. Nom va narx shu yerda "surat"
// (snapshot) qilinadi - keyinchalik Catalog'da o'zgarsa ham allaqachon berilgan buyurtmaga
// ta'sir qilmaydi.
func buildOrderItemSnapshots(ctx context.Context, productRepo catalog.ProductRepository, requests []productQuantity) ([]domain.OrderItem, error) {
	items := make([]domain.OrderItem, 0, len(requests))
	for _, req := range requests {
		product, err := productRepo.FindByID(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, catalog.ErrProductNotFound
		}

		unitPrice := product.Price()
		if discount := product.DiscountPrice(); discount != nil {
			unitPrice = *discount
		}

		item, err := domain.NewOrderItem(product.ID(), product.NameUz(), unitPrice, req.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// reserveOrderStock - har bir item uchun StockReserver orqali zaxirani atomik ravishda
// kamaytiradi. Birortasi muvaffaqiyatsiz bo'lsa, shu chaqiruv ichida kamaytirilganlarni
// avtomatik ravishda orqaga qaytaradi (hammasi yoki hech narsa) va xatoni qaytaradi.
func reserveOrderStock(ctx context.Context, stockReserver domain.StockReserver, items []domain.OrderItem) error {
	decremented := make([]domain.OrderItem, 0, len(items))
	for _, item := range items {
		if err := stockReserver.DecrementStock(ctx, item.ProductID(), item.Quantity()); err != nil {
			releaseOrderStock(ctx, stockReserver, decremented)
			return err
		}
		decremented = append(decremented, item)
	}
	return nil
}

// releaseOrderStock - berilgan itemlar uchun avval kamaytirilgan zaxirani orqaga qaytaradi
// (kompensatsiya) - masalan reserveOrderStock muvaffaqiyatli bo'lgandan keyin buyurtmani
// yaratish yoki saqlash muvaffaqiyatsiz tugasa. RestoreStock'ning o'zi xato qaytarsa ham,
// ustuvor xato allaqachon chaqiruvchida borligi uchun bu yerda e'tiborsiz qoldiriladi.
func releaseOrderStock(ctx context.Context, stockReserver domain.StockReserver, items []domain.OrderItem) {
	for _, item := range items {
		_ = stockReserver.RestoreStock(ctx, item.ProductID(), item.Quantity())
	}
}
