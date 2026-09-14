package http

import (
	"encoding/json"
	"net/http"

	identitydomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/interfaces/http/middleware"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/application"
	orderingdomain "github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	checkoutUseCase             *application.CheckoutUseCase
	getUserOrdersUseCase        *application.GetUserOrdersUseCase
	getOrderByIDUseCase         *application.GetOrderByIDUseCase
	updatePaymentStatusUseCase  *application.UpdatePaymentStatusUseCase
	updateDeliveryStatusUseCase *application.UpdateDeliveryStatusUseCase
	getAllOrdersUseCase         *application.GetAllOrdersUseCase
	createManualOrderUseCase    *application.CreateManualOrderUseCase
	cancelOrderUseCase          *application.CancelOrderUseCase
}

func NewOrderHandler(
	checkoutUC *application.CheckoutUseCase,
	getUserOrdersUC *application.GetUserOrdersUseCase,
	getOrderByIDUC *application.GetOrderByIDUseCase,
	updatePaymentStatusUC *application.UpdatePaymentStatusUseCase,
	updateDeliveryStatusUC *application.UpdateDeliveryStatusUseCase,
	getAllOrdersUC *application.GetAllOrdersUseCase,
	createManualOrderUC *application.CreateManualOrderUseCase,
	cancelOrderUC *application.CancelOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		checkoutUseCase:             checkoutUC,
		getUserOrdersUseCase:        getUserOrdersUC,
		getOrderByIDUseCase:         getOrderByIDUC,
		updatePaymentStatusUseCase:  updatePaymentStatusUC,
		updateDeliveryStatusUseCase: updateDeliveryStatusUC,
		getAllOrdersUseCase:         getAllOrdersUC,
		createManualOrderUseCase:    createManualOrderUC,
		cancelOrderUseCase:          cancelOrderUC,
	}
}

// Checkout godoc
//
//	@Summary		Savatni buyurtmaga aylantirish
//	@Description	Joriy foydalanuvchining savatidagi mahsulotlarni buyurtmaga aylantiradi: Catalog'dan joriy nom/narxni "surat" qiladi, zaxirani atomik ravishda kamaytiradi va savatni bo'shatadi
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CheckoutRequest	true	"Yetkazib berish ma'lumotlari"
//	@Success		201		{object}	response.Envelope{data=OrderResponse}	"Buyurtma yaratildi"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov, bo'sh savat yoki yaroqsiz telefon/manzil"
//	@Failure		401		{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		404		{object}	response.Envelope	"Mahsulot topilmadi"
//	@Failure		409		{object}	response.Envelope	"Mahsulot uchun yetarli zaxira yo'q"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/checkout [post]
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeOrderError(w, ErrUnauthorized)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.checkoutUseCase.Execute(r.Context(), application.CheckoutInput{
		UserID:  userID,
		Address: req.Address,
		Phone:   req.Phone,
		Note:    req.Note,
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponse(order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusCreated, resp)
}

// GetUserOrders godoc
//
//	@Summary		Mening buyurtmalarim
//	@Description	Joriy autentifikatsiya qilingan foydalanuvchining barcha buyurtmalari ro'yxati (buyurtma holati sahifasi uchun)
//	@Tags			orders
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Envelope{data=[]OrderResponse}	"Buyurtmalar ro'yxati"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/orders [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeOrderError(w, ErrUnauthorized)
		return
	}

	orders, err := h.getUserOrdersUseCase.Execute(r.Context(), userID)
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponses(orders)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusOK, resp)
}

// GetOrderByID godoc
//
//	@Summary		Buyurtma tafsilotlari
//	@Description	Bitta buyurtmani ID bo'yicha oladi. Faqat buyurtma egasi yoki admin ko'ra oladi.
//	@Tags			orders
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Buyurtma ID"
//	@Success		200	{object}	response.Envelope{data=OrderResponse}	"Buyurtma"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403	{object}	response.Envelope	"Bu buyurtmani ko'rish huquqi yo'q"
//	@Failure		404	{object}	response.Envelope	"Buyurtma topilmadi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeOrderError(w, ErrUnauthorized)
		return
	}

	// Rol context'da bo'lmasa yoki noma'lum bo'lsa, xavfsiz tomonga: admin emas deb hisoblanadi.
	role, _ := middleware.RoleFromContext(r.Context())
	isAdmin := role == identitydomain.RoleAdmin

	orderID := chi.URLParam(r, "id")

	order, err := h.getOrderByIDUseCase.Execute(r.Context(), orderID, userID, isAdmin)
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponse(order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusOK, resp)
}

// UpdatePaymentStatus godoc
//
//	@Summary		To'lov holatini o'zgartirish
//	@Description	Buyurtmaning to'lov holatini o'zgartiradi (unpaid/paid). Faqat admin uchun.
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Buyurtma ID"
//	@Param			request	body		UpdatePaymentStatusRequest	true	"Yangi to'lov holati"
//	@Success		200		{object}	response.Envelope{data=OrderResponse}	"Yangilangan buyurtma"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov yoki noto'g'ri holat qiymati"
//	@Failure		401		{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403		{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		404		{object}	response.Envelope	"Buyurtma topilmadi"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/orders/{id}/payment-status [patch]
func (h *OrderHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	var req UpdatePaymentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.updatePaymentStatusUseCase.Execute(r.Context(), orderID, orderingdomain.PaymentStatus(req.Status))
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponse(order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusOK, resp)
}

// UpdateDeliveryStatus godoc
//
//	@Summary		Yetkazib berish holatini o'zgartirish
//	@Description	Buyurtmaning yetkazib berish holatini o'zgartiradi (preparing/handed_to_courier/delivered). Faqat oldinga qarab o'zgarishi mumkin. Faqat admin uchun.
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Buyurtma ID"
//	@Param			request	body		UpdateDeliveryStatusRequest	true	"Yangi yetkazib berish holati"
//	@Success		200		{object}	response.Envelope{data=OrderResponse}	"Yangilangan buyurtma"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov yoki noto'g'ri holat qiymati"
//	@Failure		401		{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403		{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		404		{object}	response.Envelope	"Buyurtma topilmadi"
//	@Failure		409		{object}	response.Envelope	"Holatni orqaga qaytarib bo'lmaydi"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/orders/{id}/delivery-status [patch]
func (h *OrderHandler) UpdateDeliveryStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	var req UpdateDeliveryStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.updateDeliveryStatusUseCase.Execute(r.Context(), orderID, orderingdomain.DeliveryStatus(req.Status))
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponse(order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusOK, resp)
}

// GetAllOrders godoc
//
//	@Summary		Barcha buyurtmalar (admin navbati)
//	@Description	Barcha foydalanuvchilarning barcha buyurtmalari ro'yxati. Faqat admin uchun.
//	@Tags			orders
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Envelope{data=[]OrderResponse}	"Buyurtmalar ro'yxati"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403	{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/orders/admin [get]
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.getAllOrdersUseCase.Execute(r.Context())
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponses(orders)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusOK, resp)
}

// CreateManualOrder godoc
//
//	@Summary		Qo'lda buyurtma yaratish (offline savdo)
//	@Description	Admin tomonidan cart'siz, to'g'ridan-to'g'ri buyurtma yaratadi (masalan, telefon orqali qabul qilingan yoki do'kondagi offline savdo). Catalog'dan joriy nom/narxni "surat" qiladi va zaxirani atomik ravishda kamaytiradi (Checkout bilan bir xil mantiq). Faqat admin uchun.
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateManualOrderRequest	true	"Mijoz ma'lumotlari, mahsulotlar ro'yxati va boshlang'ich to'lov holati"
//	@Success		201		{object}	response.Envelope{data=OrderResponse}	"Buyurtma yaratildi"
//	@Failure		400		{object}	response.Envelope	"Noto'g'ri so'rov, bo'sh ro'yxat yoki noto'g'ri holat/manzil/telefon qiymati"
//	@Failure		401		{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403		{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		404		{object}	response.Envelope	"Mahsulot topilmadi"
//	@Failure		409		{object}	response.Envelope	"Mahsulot uchun yetarli zaxira yo'q"
//	@Failure		500		{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/orders [post]
func (h *OrderHandler) CreateManualOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateManualOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	items := make([]application.ManualOrderItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, application.ManualOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.createManualOrderUseCase.Execute(r.Context(), application.CreateManualOrderInput{
		Address:       req.Address,
		Phone:         req.Phone,
		Note:          req.Note,
		Items:         items,
		PaymentStatus: orderingdomain.PaymentStatus(req.PaymentStatus),
	})
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponse(order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusCreated, resp)
}

// CancelOrder godoc
//
//	@Summary		Buyurtmani bekor qilish
//	@Description	Buyurtmani bekor qiladi va Checkout/CreateManualOrder paytida kamaytirilgan zaxirani Catalog'ga qaytaradi. Faqat Preparing holatidagi buyurtmalar bekor qilinishi mumkin - courier'ga topshirilgan yoki yetkazilgan buyurtmani bu yo'l bilan bekor qilib bo'lmaydi. Faqat admin uchun.
//	@Tags			orders
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Buyurtma ID"
//	@Success		200	{object}	response.Envelope{data=OrderResponse}	"Bekor qilingan buyurtma"
//	@Failure		401	{object}	response.Envelope	"Autentifikatsiya talab qilinadi"
//	@Failure		403	{object}	response.Envelope	"Faqat admin uchun"
//	@Failure		404	{object}	response.Envelope	"Buyurtma topilmadi"
//	@Failure		409	{object}	response.Envelope	"Courier'ga topshirilgan yoki yetkazilgan buyurtmani bekor qilib bo'lmaydi"
//	@Failure		500	{object}	response.Envelope	"Ichki server xatosi"
//	@Router			/admin/orders/{id}/cancel [patch]
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	order, err := h.cancelOrderUseCase.Execute(r.Context(), orderID)
	if err != nil {
		writeOrderError(w, err)
		return
	}

	resp, err := ToOrderResponse(order)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.Success(w, http.StatusOK, resp)
}
