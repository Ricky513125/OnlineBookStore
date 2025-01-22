package model

type CartBook struct {
	ID              uint64
	Name, Author    string
	Quantity        uint
	Price, Subtotal float64
}

type CartBooks = []CartBook

func AddCartItem(userID, bookID uint64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := "insert ignore into cart(user_id, book_id) values(?, ?)"
	result, err := tx.Exec(query, userID, bookID)
	if err != nil {
		return err
	}

	if count, _ := result.RowsAffected(); count == 0 {
		query = `update cart
		set quantity = quantity + 1
		where user_id = ? and book_id = ?`
		tx.Exec(query, userID, bookID)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetCartItems(userID uint64) (books CartBooks) {
	query := `select id, name, author, price, quantity, price * quantity
	from book, cart
	where id = book_id
		and not deleted
		and user_id = ?`
	rs, _ := db.Query(query, userID)
	defer rs.Close()

	for book := (CartBook{}); rs.Next(); books = append(books, book) {
		rs.Scan(&book.ID, &book.Name, &book.Author, &book.Price, &book.Quantity, &book.Subtotal)
	}
	return
}

func UpdateCartItem(userID, bookID uint64, quantity uint) (err error) {
	if quantity == 0 {
		query := `delete from cart
		where user_id =? and book_id = ?`
		_, err = db.Exec(query, userID, bookID)
	} else {
		query := `update cart
		set quantity = ?
		where user_id = ? and book_id = ?`
		_, err = db.Exec(query, quantity, userID, bookID)
	}
	return
}
