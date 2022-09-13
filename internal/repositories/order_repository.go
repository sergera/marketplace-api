package repositories

import (
	"database/sql"

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

func (o *OrderRepository) UpdateOrder(m domain.OrderModel) error {
	if _, err := o.conn.Session.Exec(
		`
		UPDATE orders
		SET status = $1
		WHERE id = $2
		`,
		m.Status, m.Id,
	); err != nil {
		return err
	}

	return nil
}

func (o *OrderRepository) GetOrderRange(m domain.OrderRangeModel) ([]domain.OrderModel, error) {
	var rows *sql.Rows
	var err error
	if m.OldestFirst {
		rows, err = o.conn.Session.Query(
			`
			SELECT id, price, status, date_created
			FROM orders
			WHERE id >= $1 AND id <= $2
			ORDER BY id ASC
			`,
			m.Start, m.End,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	} else {
		rows, err = o.conn.Session.Query(
			`
			WITH last_order AS (
				SELECT id
				FROM orders
				ORDER BY id DESC
				LIMIT 1
			)

			SELECT id, price, status, date_created
			FROM orders
			WHERE id <= ((SELECT id from last_order) - ($1 - 1))
			AND id >= ((SELECT id from last_order) - ($2 - 1))
			ORDER BY id DESC
			`,
			m.Start, m.End,
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

	return orders, nil
}
