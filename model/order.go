package model

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

type OrderItem struct {
	Name     string
	Cover    sql.NullString
	Quantity uint
	Price    float64
	Subtotal float64
}

type OrderItems = []OrderItem

type Order struct {
	ID        uint64
	UserID    uint64
	AddressID uint64
	Total     float64
	Status    string
	Items     OrderItems
	time.Time
}

type Orders = []Order

func AddOrder(userID uint64, addressID uint, booksID []uint64) (uint64, error) {
	if len(booksID) == 0 {
		return 0, nil
	}

	txn, _ := db.Begin()

	query := `insert into orders(total, user_id, address_id, status)
	values(0, ?, ?, ?)`
	result, err := txn.Exec(query, userID, addressID, "unpaid")
	if err != nil {
		return 0, err
	}
	orderID, _ := result.LastInsertId()

	holders := strings.Repeat(",?", len(booksID))[1:]
	query = `insert into order_item(order_id, book_id, price, quantity)
	select ?, book_id, price, quantity
	from book, cart
	where id = book_id and user_id = ? and book_id in (` + holders + ")"

	args := make([]interface{}, len(booksID)+2)
	args[0], args[1] = orderID, userID
	for i, v := range booksID {
		args[i+2] = v
	}
	_, err = txn.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	query = `update orders
	set total = (select sum(subtotal)
			from order_item
			where order_id = ?)
	where id = ?`
	_, err = txn.Exec(query, orderID, orderID)
	if err != nil {
		return 0, err
	}

	query = `delete from cart
	where user_id = ? and book_id in (` + holders + ")"
	_, err = txn.Exec(query, args[1:]...)
	if err != nil {
		return 0, err
	}

	txn.Commit()
	return uint64(orderID), nil
}

func CountOrders(userID uint64) (count uint64, err error) {
	query := `select count(*)
	from orders`

	args := []any{}
	if userID != 0 {
		query += ` where user_id = ?`
		args = append(args, userID)
	}

	err = db.QueryRow(query, args...).Scan(&count)
	return
}

func GetOrders(userID, offset, row_count uint64) (orders Orders, err error) {
	query := `select id, order_time, total, status
	from orders`
	args := []any{}

	if userID != 0 {
		query += " where user_id = ? "
		args = append(args, userID)
	}

	query += " order by id desc "

	var rs *sql.Rows
	if row_count != 0 {
		query += "limit ?, ?"
		args = append(args, offset, row_count)
	}
	rs, err = db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	var order Order
	for rs.Next() {
		err := rs.Scan(&order.ID, &order.Time, &order.Total, &order.Status)
		if err != nil {
			log.Println(err)
		} else {
			order.Items, _ = GetOrderItems(order.ID)
			orders = append(orders, order)
		}
	}
	return
}

func GetOrder(id uint64, needItems bool) (order Order, err error) {
	order.ID = id

	query := `select user_id, order_time, total, status, address_id
	from orders
	where id = ?`
	err = db.QueryRow(query, id).
		Scan(&order.UserID, &order.Time, &order.Total, &order.Status, &order.AddressID)
	if err != nil {
		return
	}

	if needItems {
		order.Items, err = GetOrderItems(id)
	}
	return
}

func UpdateOrderStatus(id uint64, status string) (err error) {
	query := `update orders
	set status = ?
	where id = ?`
	_, err = db.Exec(query, status, id)

	// 增加成功购买量
	if status == "success" {
		err := AddPurchasedNum(id)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

func GetOrderItems(orderID uint64) (orderItems OrderItems, err error) {
	query := `select name, cover, quantity, order_item.price, subtotal
	from order_item, book
	where book_id = book.id and order_id = ?`
	rs, err := db.Query(query, orderID)
	if err != nil {
		return nil, err
	}

	var orderItem OrderItem
	for rs.Next() {
		err := rs.Scan(&orderItem.Name, &orderItem.Cover, &orderItem.Quantity, &orderItem.Price, &orderItem.Subtotal)
		if err != nil {
			log.Println(err)
		} else {
			orderItems = append(orderItems, orderItem)
		}
	}
	return
}
