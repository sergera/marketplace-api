package repositories

import (
	"database/sql"
	"strconv"

	"github.com/sergera/marketplace-api/internal/domain"
)

type OrderRepository struct {
	conn *DBConnection
}

func NewOrderRepository(conn *DBConnection) *OrderRepository {
	return &OrderRepository{conn}
}

func (o *OrderRepository) CreateOrder(m *domain.OrderModel) error {
	if err := o.conn.Session.QueryRow(
		`
		INSERT INTO orders (price, status, date_created)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		m.Price, m.Status, m.Date,
	).Scan(&m.Id); err != nil {
		return err
	}

	return nil
}

func (o *OrderRepository) GetOrderRange(m domain.OrderRangeModel) ([]domain.OrderModel, error) {
	tx, err := o.conn.Session.Begin()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if m.OldestFirst {
		rows, err = tx.Query(
			`
			SELECT id, price, status, date_created
			FROM orders
			AND id >= $1
			AND id <= $2
			ORDER BY id ASC
			`,
			m.Start, m.End,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	} else {
		var maxIdString string
		tx.QueryRow(
			`
			SELECT id
			FROM orders
			ORDER BY id DESC
			LIMIT 1
			`,
		).Scan(&maxIdString)

		if maxIdString == "" {
			var emptyOrderSlice []domain.OrderModel
			return emptyOrderSlice, nil
		}

		maxId, err := strconv.ParseInt(maxIdString, 10, 64)
		if err != nil {
			return nil, err
		}

		start, err := strconv.ParseInt(m.Start, 10, 64)
		if err != nil {
			return nil, err
		}

		end, err := strconv.ParseInt(m.End, 10, 64)
		if err != nil {
			return nil, err
		}

		rows, err = tx.Query(
			`
			SELECT id, price, status, date_created
			FROM orders
			AND id >= $1
			AND id <= $2
			ORDER BY id DESC
			`,
			maxId-(end-1), maxId-(start-1),
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	}

	var orders []domain.OrderModel
	for rows.Next() {
		order := domain.OrderModel{}
		err := rows.Scan(&order.Id, &order.Price, &order.Status, &order.Date)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return orders, nil
}
