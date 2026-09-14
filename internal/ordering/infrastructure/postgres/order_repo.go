package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/ordering/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/money"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresOrderRepository struct {
	db *sqlx.DB
}

func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

const orderColumns = `id, user_id, customer_address, customer_phone, customer_note,
	payment_status, delivery_status, created_at`

const orderItemColumns = `order_id, product_id, product_name, unit_price_amount, unit_price_currency, quantity`

type orderRow struct {
	ID              uuid.UUID      `db:"id"`
	UserID          sql.NullString `db:"user_id"`
	CustomerAddress string         `db:"customer_address"`
	CustomerPhone   string         `db:"customer_phone"`
	CustomerNote    sql.NullString `db:"customer_note"`
	PaymentStatus   string         `db:"payment_status"`
	DeliveryStatus  string         `db:"delivery_status"`
	CreatedAt       time.Time      `db:"created_at"`
}

type orderItemRow struct {
	OrderID           uuid.UUID `db:"order_id"`
	ProductID         string    `db:"product_id"`
	ProductName       string    `db:"product_name"`
	UnitPriceAmount   float64   `db:"unit_price_amount"`
	UnitPriceCurrency string    `db:"unit_price_currency"`
	Quantity          int       `db:"quantity"`
}

// Save - yangi buyurtmani order va order_items jadvallariga bitta tranzaksiyada yozadi.
func (r *PostgresOrderRepository) Save(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Commit muvaffaqiyatli bo'lsa, bu no-op bo'ladi

	info := order.CustomerInfo()

	orderQuery := `INSERT INTO orders (` + orderColumns + `) VALUES (
		:id, :user_id, :customer_address, :customer_phone, :customer_note,
		:payment_status, :delivery_status, :created_at
	)`
	orderParams := map[string]any{
		"id":               order.ID(),
		"user_id":          nullableString(order.UserID()),
		"customer_address": info.Address(),
		"customer_phone":   info.Phone(),
		"customer_note":    nullIfEmpty(info.Note()),
		"payment_status":   string(order.PaymentStatus()),
		"delivery_status":  string(order.DeliveryStatus()),
		"created_at":       order.CreatedAt(),
	}
	if _, err := tx.NamedExecContext(ctx, orderQuery, orderParams); err != nil {
		return err
	}

	itemQuery := `INSERT INTO order_items (` + orderItemColumns + `) VALUES (
		:order_id, :product_id, :product_name, :unit_price_amount, :unit_price_currency, :quantity
	)`
	for _, item := range order.Items() {
		itemParams := map[string]any{
			"order_id":            order.ID(),
			"product_id":          item.ProductID(),
			"product_name":        item.ProductName(),
			"unit_price_amount":   item.UnitPrice().Amount(),
			"unit_price_currency": item.UnitPrice().Currency(),
			"quantity":            item.Quantity(),
		}
		if _, err := tx.NamedExecContext(ctx, itemQuery, itemParams); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders WHERE id = :id`
	params := map[string]any{"id": id}

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var row orderRow
	if err := stmt.GetContext(ctx, &row, params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // topilmadi — bu xato emas, use case ErrOrderNotFound'ga aylantiradi
		}
		return nil, err
	}

	items, err := r.findItems(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return rowToOrder(row, items), nil
}

func (r *PostgresOrderRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders WHERE user_id = :user_id ORDER BY created_at DESC`
	params := map[string]any{"user_id": userID}

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []orderRow
	if err := stmt.SelectContext(ctx, &rows, params); err != nil {
		return nil, err
	}

	return r.rowsToOrders(ctx, rows)
}

// FindAll - admin uchun barcha buyurtmalar ro'yxati. Bu so'rovda hech qanday parametr yo'q,
// shuning uchun named yoki positional bog'lash shart emas.
func (r *PostgresOrderRepository) FindAll(ctx context.Context) ([]*domain.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders ORDER BY created_at DESC`

	var rows []orderRow
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, err
	}

	return r.rowsToOrders(ctx, rows)
}

// Update - Order aggregate'ining faqat o'zgarishi mumkin bo'lgan maydonlarini (to'lov va
// yetkazib berish holatlari) yangilaydi. CustomerInfo, Items va CreatedAt yaratilgandan keyin
// o'zgarmaydi, shuning uchun bu yerda qayta yozilmaydi.
func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `UPDATE orders SET payment_status = :payment_status, delivery_status = :delivery_status WHERE id = :id`
	params := map[string]any{
		"payment_status":  string(order.PaymentStatus()),
		"delivery_status": string(order.DeliveryStatus()),
		"id":              order.ID(),
	}

	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (r *PostgresOrderRepository) findItems(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	query := `SELECT ` + orderItemColumns + ` FROM order_items WHERE order_id = :order_id`
	params := map[string]any{"order_id": orderID}

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []orderItemRow
	if err := stmt.SelectContext(ctx, &rows, params); err != nil {
		return nil, err
	}

	items := make([]domain.OrderItem, 0, len(rows))
	for _, row := range rows {
		price, err := money.NewMoney(row.UnitPriceAmount, row.UnitPriceCurrency)
		if err != nil {
			return nil, err
		}
		items = append(items, domain.NewOrderItemFromRepository(row.ProductID, row.ProductName, price, row.Quantity))
	}
	return items, nil
}

// rowsToOrders - har bir order qatori uchun order_items'ni alohida so'rov bilan biriktiradi.
// N+1 - bu qatlamda repository'lar bo'ylab umumiy JOIN/batch strategiyasi hali yo'q (boshqa
// context'larda ham xuddi shunday), lekin har bir so'rov named parametr bilan qoladi.
func (r *PostgresOrderRepository) rowsToOrders(ctx context.Context, rows []orderRow) ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0, len(rows))
	for _, row := range rows {
		items, err := r.findItems(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, rowToOrder(row, items))
	}
	return orders, nil
}

func rowToOrder(row orderRow, items []domain.OrderItem) *domain.Order {
	var userID *string
	if row.UserID.Valid {
		v := row.UserID.String
		userID = &v
	}

	note := ""
	if row.CustomerNote.Valid {
		note = row.CustomerNote.String
	}

	info := domain.NewCustomerInfoFromRepository(row.CustomerAddress, row.CustomerPhone, note)

	return domain.NewOrderFromRepository(
		row.ID,
		userID,
		info,
		items,
		domain.PaymentStatus(row.PaymentStatus),
		domain.DeliveryStatus(row.DeliveryStatus),
		row.CreatedAt,
	)
}

func nullableString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func nullIfEmpty(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
