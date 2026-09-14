package postgres

import (
	"context"
	"errors"

	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// CatalogStockReserver - domain.StockReserver ning implementatsiyasi. Catalog'ning
// ProductRepository'siga to'g'ridan-to'g'ri murojaat qiladi - bu cart allaqachon catalog'ga
// nisbatan qiladigan xil infra-dan-infra-ga bog'lanish (bu chegara hali tuzatilmagan, shu
// bilan birga izchil davom ettirilmoqda). Catalog'ning ErrInsufficientStock xatosi ordering'ning
// o'z sentinel xatosiga tarjima qilinadi - chaqiruvchi qatlamlar (CheckoutUseCase, keyinchalik
// HTTP xato xaritalash) faqat ordering domenining o'z xato lug'atiga tayanadi.
type CatalogStockReserver struct {
	productRepo catalog.ProductRepository
}

func NewCatalogStockReserver(productRepo catalog.ProductRepository) *CatalogStockReserver {
	return &CatalogStockReserver{productRepo: productRepo}
}

func (r *CatalogStockReserver) DecrementStock(ctx context.Context, productID string, quantity int) error {
	err := r.productRepo.DecrementStock(ctx, productID, quantity)
	if errors.Is(err, catalog.ErrInsufficientStock) {
		return domain.ErrInsufficientStock
	}
	return err
}

func (r *CatalogStockReserver) RestoreStock(ctx context.Context, productID string, quantity int) error {
	return r.productRepo.IncrementStock(ctx, productID, quantity)
}
