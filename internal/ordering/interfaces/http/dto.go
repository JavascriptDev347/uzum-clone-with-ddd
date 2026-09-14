package http

import (
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
)

// CheckoutRequest - checkout so'rovining tanasi. Note ixtiyoriy.
type CheckoutRequest struct {
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Note    string `json:"note,omitempty"`
}

// UpdatePaymentStatusRequest - admin uchun to'lov holatini o'zgartirish so'rovi. Status "unpaid"
// yoki "paid" bo'lishi kerak.
type UpdatePaymentStatusRequest struct {
	Status string `json:"status"`
}

// UpdateDeliveryStatusRequest - admin uchun yetkazib berish holatini o'zgartirish so'rovi.
// Status "preparing", "handed_to_courier" yoki "delivered" bo'lishi kerak, va faqat oldinga
// qarab o'zgarishi mumkin.
type UpdateDeliveryStatusRequest struct {
	Status string `json:"status"`
}

// ManualOrderItemRequest - admin tomonidan qo'lda kiritiladigan bitta qator: mahsulot va miqdor.
type ManualOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// CreateManualOrderRequest - admin uchun cart'siz, to'g'ridan-to'g'ri buyurtma yaratish so'rovi
// (offline savdo). PaymentStatus "unpaid" yoki "paid" bo'lishi kerak - offline savdolar
// ko'pincha joyida to'lanadi, shuning uchun standart Unpaid'dan farqli qiymat berilishi mumkin.
type CreateManualOrderRequest struct {
	Address       string                   `json:"address"`
	Phone         string                   `json:"phone"`
	Note          string                   `json:"note,omitempty"`
	Items         []ManualOrderItemRequest `json:"items"`
	PaymentStatus string                   `json:"payment_status"`
}

// OrderItemResponse - buyurtma qatorining HTTP javob shakli. domain.OrderItem to'g'ridan-to'g'ri
// serializatsiya qilinmaydi (private field'lari bor), shu sabab bu alohida DTO orqali map qilinadi.
type OrderItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Currency    string  `json:"currency"`
	Quantity    int     `json:"quantity"`
}

// OrderResponse - buyurtmaning HTTP javob shakli. domain.Order va domain.CustomerInfo
// to'g'ridan-to'g'ri serializatsiya qilinmaydi - bu DTO orqali map qilinadi.
type OrderResponse struct {
	ID             string              `json:"id"`
	UserID         *string             `json:"user_id,omitempty"`
	Address        string              `json:"address"`
	Phone          string              `json:"phone"`
	Note           string              `json:"note,omitempty"`
	Items          []OrderItemResponse `json:"items"`
	PaymentStatus  string              `json:"payment_status"`
	DeliveryStatus string              `json:"delivery_status"`
	TotalAmount    float64             `json:"total_amount"`
	TotalCurrency  string              `json:"total_currency"`
	CreatedAt      time.Time           `json:"created_at"`
}

// ToOrderResponse - domain.Order'ni HTTP javob DTO'siga aylantiradi. Order.Total() qatorlar
// valyutasi mos kelmasa xato qaytarishi mumkin (amalda checkout doim bitta valyutadan
// foydalanadi), shuning uchun bu yerda ham xato tarqatiladi.
func ToOrderResponse(order *domain.Order) (OrderResponse, error) {
	items := make([]OrderItemResponse, 0, len(order.Items()))
	for _, item := range order.Items() {
		items = append(items, OrderItemResponse{
			ProductID:   item.ProductID(),
			ProductName: item.ProductName(),
			UnitPrice:   item.UnitPrice().Amount(),
			Currency:    item.UnitPrice().Currency(),
			Quantity:    item.Quantity(),
		})
	}

	total, err := order.Total()
	if err != nil {
		return OrderResponse{}, err
	}

	info := order.CustomerInfo()

	return OrderResponse{
		ID:             order.ID().String(),
		UserID:         order.UserID(),
		Address:        info.Address(),
		Phone:          info.Phone(),
		Note:           info.Note(),
		Items:          items,
		PaymentStatus:  string(order.PaymentStatus()),
		DeliveryStatus: string(order.DeliveryStatus()),
		TotalAmount:    total.Amount(),
		TotalCurrency:  total.Currency(),
		CreatedAt:      order.CreatedAt(),
	}, nil
}

func ToOrderResponses(orders []*domain.Order) ([]OrderResponse, error) {
	responses := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		resp, err := ToOrderResponse(o)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}
	return responses, nil
}
