package domain

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/money"
	"github.com/google/uuid"
)

// PaymentStatus - buyurtmaning to'lov holati.
type PaymentStatus string

const (
	PaymentStatusUnpaid PaymentStatus = "unpaid"
	PaymentStatusPaid   PaymentStatus = "paid"
)

func (s PaymentStatus) IsValid() bool {
	return s == PaymentStatusUnpaid || s == PaymentStatusPaid
}

// DeliveryStatus - buyurtmaning yetkazib berish holati. Holatlar faqat oldinga qarab
// o'zgarishi mumkin (Preparing -> HandedToCourier -> Delivered), orqaga qaytarib bo'lmaydi.
type DeliveryStatus string

const (
	DeliveryStatusPreparing       DeliveryStatus = "preparing"
	DeliveryStatusHandedToCourier DeliveryStatus = "handed_to_courier"
	DeliveryStatusDelivered       DeliveryStatus = "delivered"

	// DeliveryStatusCancelled - alohida, chiziqli progressiyadan tashqaridagi terminal holat.
	// Faqat Preparing'dan (Cancel() orqali) erishiladi va undan boshqa hech qanday holatga
	// o'tib bo'lmaydi - shuning uchun bu deliveryRank'ga qo'shilmagan va AdvanceDelivery
	// orqali umuman erishib bo'lmaydi (pastga qarang).
	DeliveryStatusCancelled DeliveryStatus = "cancelled"
)

func (s DeliveryStatus) IsValid() bool {
	switch s {
	case DeliveryStatusPreparing, DeliveryStatusHandedToCourier, DeliveryStatusDelivered, DeliveryStatusCancelled:
		return true
	}
	return false
}

// deliveryRank - DeliveryStatus'lar orasidagi tartibni belgilaydi, orqaga qaytishni aniqlash
// uchun. Faqat oldinga siljiydigan uchta holatni o'z ichiga oladi - Cancelled bu yerga
// qo'shilmaydi, chunki u chiziqli emas (faqat Preparing'dan, faqat bitta yo'nalishda).
var deliveryRank = map[DeliveryStatus]int{
	DeliveryStatusPreparing:       0,
	DeliveryStatusHandedToCourier: 1,
	DeliveryStatusDelivered:       2,
}

// OrderItem - buyurtma qatori. ProductName va UnitPrice buyurtma yaratilgan paytdagi
// qiymatlarning suratidir (snapshot) - keyinchalik catalog'dagi mahsulot nomi yoki narxi
// o'zgarsa ham, allaqachon berilgan buyurtma o'zgarmasligi kerak.
type OrderItem struct {
	productID   string
	productName string
	unitPrice   money.Money
	quantity    int
}

func NewOrderItem(productID, productName string, unitPrice money.Money, quantity int) (OrderItem, error) {
	if productID == "" {
		return OrderItem{}, ErrEmptyOrderItemProductID
	}
	if productName == "" {
		return OrderItem{}, ErrEmptyOrderItemName
	}
	if quantity <= 0 {
		return OrderItem{}, ErrInvalidQuantity
	}

	return OrderItem{
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}, nil
}

// NewOrderItemFromRepository - saqlangan OrderItem'ni bazadan qayta tiklash uchun,
// validatsiyadan ajratilgan (invariantlar bazadan o'qishda qayta tekshirilmaydi).
func NewOrderItemFromRepository(productID, productName string, unitPrice money.Money, quantity int) OrderItem {
	return OrderItem{
		productID:   productID,
		productName: productName,
		unitPrice:   unitPrice,
		quantity:    quantity,
	}
}

func (i OrderItem) ProductID() string      { return i.productID }
func (i OrderItem) ProductName() string    { return i.productName }
func (i OrderItem) UnitPrice() money.Money { return i.unitPrice }
func (i OrderItem) Quantity() int          { return i.quantity }

// Subtotal - shu qatorning umumiy summasi (UnitPrice * Quantity).
func (i OrderItem) Subtotal() (money.Money, error) {
	return money.NewMoney(i.unitPrice.Amount()*float64(i.quantity), i.unitPrice.Currency())
}

// Order - aggregate root. UserID nil bo'lishi mumkin - bu ro'yxatdan o'tmagan mijoz uchun
// admin tomonidan qo'lda (offline) yaratilgan buyurtmani bildiradi.
type Order struct {
	id             uuid.UUID
	userID         *string
	customerInfo   CustomerInfo
	items          []OrderItem
	paymentStatus  PaymentStatus
	deliveryStatus DeliveryStatus
	createdAt      time.Time
}

// NewOrder - yangi buyurtma yaratish uchun. PaymentStatus va DeliveryStatus mos ravishda
// Unpaid va Preparing bilan boshlanadi.
func NewOrder(userID *string, customerInfo CustomerInfo, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrderItems
	}

	return &Order{
		id:             uuid.New(),
		userID:         userID,
		customerInfo:   customerInfo,
		items:          items,
		paymentStatus:  PaymentStatusUnpaid,
		deliveryStatus: DeliveryStatusPreparing,
		createdAt:      time.Now(),
	}, nil
}

// NewOrderFromRepository - saqlangan buyurtmani bazadan qayta tiklash uchun, domen
// konstruktoridan ajratilgan (invariantlar bazadan o'qishda qayta tekshirilmaydi).
func NewOrderFromRepository(
	id uuid.UUID,
	userID *string,
	customerInfo CustomerInfo,
	items []OrderItem,
	paymentStatus PaymentStatus,
	deliveryStatus DeliveryStatus,
	createdAt time.Time,
) *Order {
	return &Order{
		id:             id,
		userID:         userID,
		customerInfo:   customerInfo,
		items:          items,
		paymentStatus:  paymentStatus,
		deliveryStatus: deliveryStatus,
		createdAt:      createdAt,
	}
}

func (o *Order) ID() uuid.UUID                  { return o.id }
func (o *Order) UserID() *string                { return o.userID }
func (o *Order) CustomerInfo() CustomerInfo     { return o.customerInfo }
func (o *Order) Items() []OrderItem             { return o.items }
func (o *Order) PaymentStatus() PaymentStatus   { return o.paymentStatus }
func (o *Order) DeliveryStatus() DeliveryStatus { return o.deliveryStatus }
func (o *Order) CreatedAt() time.Time           { return o.createdAt }

func (o *Order) MarkPaid() {
	o.paymentStatus = PaymentStatusPaid
}

func (o *Order) MarkUnpaid() {
	o.paymentStatus = PaymentStatusUnpaid
}

// AdvanceDelivery - yetkazib berish holatini o'zgartiradi. Faqat joriy holatdan oldinga
// (yoki bir xil holatga) o'tishga ruxsat beradi, orqaga qaytarishga urinilsa xatolik qaytaradi.
// Cancelled bu yo'l orqali umuman o'rnatilmaydi (buyurtmani bekor qilish uchun Cancel()
// ishlatilishi kerak) va allaqachon Cancelled bo'lgan buyurtma bu metod orqali boshqa hech
// qanday holatga o'tkazilmaydi - u terminal holat.
func (o *Order) AdvanceDelivery(status DeliveryStatus) error {
	if !status.IsValid() {
		return ErrInvalidDeliveryStatus
	}
	if o.deliveryStatus == DeliveryStatusCancelled || status == DeliveryStatusCancelled {
		return ErrInvalidDeliveryTransition
	}
	if deliveryRank[status] < deliveryRank[o.deliveryStatus] {
		return ErrInvalidDeliveryTransition
	}
	o.deliveryStatus = status
	return nil
}

// Cancel - buyurtmani bekor qiladi. Faqat Preparing holatidagi buyurtmalar bekor qilinishi
// mumkin - courier'ga allaqachon topshirilgan yoki yetkazilgan buyurtmani bu yo'l bilan
// jimgina "bekor qilib" bo'lmaydi. Bu bir yo'nalishli, terminal holat o'tishi bo'lgani uchun
// AdvanceDelivery'ning umumiy rank-asosli mantig'iga to'g'ri kelmaydi, shuning uchun alohida
// metod sifatida ajratilgan.
func (o *Order) Cancel() error {
	if o.deliveryStatus != DeliveryStatusPreparing {
		return ErrOrderAlreadyShipped
	}
	o.deliveryStatus = DeliveryStatusCancelled
	return nil
}

// Total - buyurtmadagi barcha qatorlar summasining yig'indisi.
func (o *Order) Total() (money.Money, error) {
	if len(o.items) == 0 {
		return money.Money{}, ErrEmptyOrderItems
	}

	total, err := o.items[0].Subtotal()
	if err != nil {
		return money.Money{}, err
	}

	for _, item := range o.items[1:] {
		subtotal, err := item.Subtotal()
		if err != nil {
			return money.Money{}, err
		}
		total, err = total.Add(subtotal)
		if err != nil {
			return money.Money{}, err
		}
	}

	return total, nil
}
