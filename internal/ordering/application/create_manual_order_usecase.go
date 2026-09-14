package application

import (
	"context"

	catalog "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// CreateManualOrderUseCase - admin uchun: cart'siz, to'g'ridan-to'g'ri buyurtma yaratish
// (masalan, telefon orqali qabul qilingan yoki do'kondagi offline savdo). UserID har doim
// nil - bu ro'yxatdan o'tmagan mijoz uchun admin tomonidan qo'lda yaratilgan buyurtmani
// bildiradi (domain.Order'ning o'zida ko'zda tutilgan holat).
type CreateManualOrderUseCase struct {
	orderRepo     domain.OrderRepository
	productRepo   catalog.ProductRepository
	stockReserver domain.StockReserver
}

func NewCreateManualOrderUseCase(
	orderRepo domain.OrderRepository,
	productRepo catalog.ProductRepository,
	stockReserver domain.StockReserver,
) *CreateManualOrderUseCase {
	return &CreateManualOrderUseCase{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		stockReserver: stockReserver,
	}
}

// ManualOrderItemInput - admin tomonidan qo'lda kiritiladigan bitta qator: mahsulot va miqdor.
type ManualOrderItemInput struct {
	ProductID string
	Quantity  int
}

type CreateManualOrderInput struct {
	Address string
	Phone   string
	Note    string
	Items   []ManualOrderItemInput

	// PaymentStatus - admin buyurtmani yaratishda boshlang'ich to'lov holatini belgilaydi
	// (offline savdolar ko'pincha joyida to'lanadi) - domain.NewOrder()'ning standart
	// Unpaid holatidan farqli o'laroq.
	PaymentStatus domain.PaymentStatus
}

// mergeDuplicateManualOrderItems - admin bir xil ProductID'ni ro'yxatda bir necha marta
// kiritishi mumkin (masalan telefon orqali qabul qilingan buyurtmani qo'lda yig'ishda).
// order_items jadvalida (order_id, product_id) composite PK bo'lgani uchun bitta
// ProductID uchun ikkita qator saqlashga urinish saqlashda constraint xatosiga olib
// keladi - shuning uchun bir xil ProductID'lar birlashtirilib, Quantity'lari qo'shiladi,
// birinchi uchragan tartib saqlangan holda.
func mergeDuplicateManualOrderItems(items []ManualOrderItemInput) []productQuantity {
	merged := make([]productQuantity, 0, len(items))
	indexByProductID := make(map[string]int, len(items))

	for _, item := range items {
		if idx, ok := indexByProductID[item.ProductID]; ok {
			merged[idx].Quantity += item.Quantity
			continue
		}
		indexByProductID[item.ProductID] = len(merged)
		merged = append(merged, productQuantity{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	return merged
}

func (uc *CreateManualOrderUseCase) Execute(ctx context.Context, in CreateManualOrderInput) (*domain.Order, error) {
	customerInfo, err := domain.NewCustomerInfo(in.Address, in.Phone, in.Note)
	if err != nil {
		return nil, err
	}

	if len(in.Items) == 0 {
		return nil, domain.ErrEmptyOrderItems
	}

	if !in.PaymentStatus.IsValid() {
		return nil, domain.ErrInvalidPaymentStatus
	}

	requests := mergeDuplicateManualOrderItems(in.Items)

	// Checkout bilan bir xil: Catalog'dan joriy nom/narxni "surat" qilamiz - keyinchalik
	// Catalog'da o'zgarsa ham bu buyurtmaga ta'sir qilmaydi.
	items, err := buildOrderItemSnapshots(ctx, uc.productRepo, requests)
	if err != nil {
		return nil, err
	}

	// Checkout bilan bir xil atomik zaxira kamaytirish (StockReserver porti orqali,
	// reserveOrderStock funksiyasi orqali baham ko'riladi - mantiq takrorlanmaydi).
	if err := reserveOrderStock(ctx, uc.stockReserver, items); err != nil {
		return nil, err
	}

	order, err := domain.NewOrder(nil, customerInfo, items)
	if err != nil {
		releaseOrderStock(ctx, uc.stockReserver, items)
		return nil, err
	}

	// Admin belgilagan boshlang'ich to'lov holatini qo'llaymiz. NewOrder() standart
	// ravishda Unpaid bilan boshlaydi, shuning uchun Paid so'ralgan taqdirdagina o'zgartiramiz.
	if in.PaymentStatus == domain.PaymentStatusPaid {
		order.MarkPaid()
	}

	if err := uc.orderRepo.Save(ctx, order); err != nil {
		releaseOrderStock(ctx, uc.stockReserver, items)
		return nil, err
	}

	return order, nil
}
