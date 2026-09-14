package postgres

import (
	"context"

	orderingapp "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
)

// OrderingPurchaseChecker - domain.DeliveredPurchaseChecker ning implementatsiyasi. Review
// ordering'ning domain yoki repository (DB) qatlamiga to'g'ridan-to'g'ri kirmaydi - faqat
// ordering'ning application qatlamidagi fokuslangan HasDeliveredProductUseCase orqali
// so'raydi. Bu ACL chegarasi review context'i uchun qat'iy saqlanadi (cart'ning catalog
// bilan bog'lanishidan farqli o'laroq, u hali infra-dan-infra-ga to'g'ridan-to'g'ri murojaat
// qiladi).
type OrderingPurchaseChecker struct {
	hasDeliveredProductUC *orderingapp.HasDeliveredProductUseCase
}

func NewOrderingPurchaseChecker(hasDeliveredProductUC *orderingapp.HasDeliveredProductUseCase) *OrderingPurchaseChecker {
	return &OrderingPurchaseChecker{hasDeliveredProductUC: hasDeliveredProductUC}
}

func (c *OrderingPurchaseChecker) HasDeliveredPurchase(ctx context.Context, userID, productID string) (string, bool, error) {
	return c.hasDeliveredProductUC.Execute(ctx, userID, productID)
}
